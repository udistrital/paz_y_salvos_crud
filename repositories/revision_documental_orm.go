package repositories

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/services"
	"github.com/udistrital/paz_y_salvos_crud/utils"
)

type formularioGradoEntity struct {
	Id                int        `orm:"column(id);pk;auto"`
	SolicitudGradoId  int        `orm:"column(solicitud_grado_id)"`
	Version           int16      `orm:"column(version)"`
	Contenido         string     `orm:"column(contenido);type(text)"`
	Activo            bool       `orm:"column(activo)"`
	FechaCreacion     time.Time  `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time  `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
	FechaRadicacion   *time.Time `orm:"column(fecha_radicacion);null;type(timestamp without time zone)"`
}

func (formularioGradoEntity) TableName() string { return "formulario_solicitud_grado" }

func init() {
	orm.RegisterModel(
		new(models.SolicitudGrado), new(formularioGradoEntity), new(models.HistorialSolicitudGrado),
		new(models.SoporteGrado), new(models.HistorialSoporteGrado), new(models.PazSalvoGrado), new(models.HistorialPazSalvoGrado),
	)
}

type RevisionDocumentalORM struct{}

func NuevaRevisionDocumentalORM() *RevisionDocumentalORM { return &RevisionDocumentalORM{} }

func estadoPermitidoRevision(estado int, permitidos []int) bool {
	for _, permitido := range permitidos {
		if estado == permitido {
			return true
		}
	}
	return false
}

func cargarSolicitudRevision(o orm.QueryExecutor, solicitud models.SolicitudGrado, estadosPermitidos []int) (*models.BorradorInscripcionGrado, error) {
	resultado := &models.BorradorInscripcionGrado{Solicitud: solicitud, Soportes: []models.SoporteGrado{}}
	if err := o.QueryTable(new(models.HistorialSolicitudGrado)).
		Filter("SolicitudGradoId", solicitud.Id).Filter("Activo", true).
		OrderBy("-Id").One(&resultado.Historial); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, services.ErrBorradorCerrado
		}
		return nil, err
	}
	if !estadoPermitidoRevision(resultado.Historial.EstadoSolicitudId, estadosPermitidos) || resultado.Historial.FormularioSolicitudGradoId == nil {
		return nil, services.ErrBorradorCerrado
	}
	var formulario formularioGradoEntity
	if err := o.QueryTable(new(formularioGradoEntity)).
		Filter("Id", *resultado.Historial.FormularioSolicitudGradoId).
		Filter("SolicitudGradoId", solicitud.Id).Filter("Activo", true).One(&formulario); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, services.ErrBorradorCerrado
		}
		return nil, err
	}
	if !json.Valid([]byte(formulario.Contenido)) {
		return nil, services.ErrBorradorCerrado
	}
	resultado.Formulario = models.FormularioSolicitudGrado{
		Id: formulario.Id, SolicitudGradoId: formulario.SolicitudGradoId, Version: formulario.Version,
		Contenido: json.RawMessage(formulario.Contenido), Activo: formulario.Activo,
		FechaCreacion: formulario.FechaCreacion, FechaModificacion: formulario.FechaModificacion,
		FechaRadicacion: formulario.FechaRadicacion,
	}
	if _, err := o.QueryTable(new(models.SoporteGrado)).
		Filter("SolicitudGradoId", solicitud.Id).
		Filter("FormularioSolicitudGradoId", resultado.Formulario.Id).
		Filter("Activo", true).OrderBy("TipoDocumentoId").All(&resultado.Soportes); err != nil {
		return nil, err
	}
	return resultado, nil
}

func (r *RevisionDocumentalORM) Consultar(id, dependenciaID int, estadosPermitidos []int) (*models.BorradorInscripcionGrado, error) {
	o := orm.NewOrm()
	var solicitud models.SolicitudGrado
	if err := o.QueryTable(new(models.SolicitudGrado)).Filter("Id", id).
		Filter("DependenciaOikosId", dependenciaID).Filter("Activo", true).One(&solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, services.ErrSolicitudFueraAlcance
		}
		return nil, err
	}
	return cargarSolicitudRevision(o, solicitud, estadosPermitidos)
}

func (r *RevisionDocumentalORM) Listar(dependenciaID int, estadosPermitidos []int, limit, offset int, filtros ...int) ([]models.BorradorInscripcionGrado, error) {
	o := orm.NewOrm()
	var solicitudes []models.SolicitudGrado
	consulta := o.QueryTable(new(models.SolicitudGrado)).Filter("DependenciaOikosId", dependenciaID).Filter("Activo", true)
	if len(filtros) > 0 && filtros[0] > 0 {
		consulta = consulta.Filter("PeriodoId", filtros[0])
	}
	if len(filtros) > 1 && filtros[1] > 0 {
		consulta = consulta.Filter("ProgramaAcademicoId", filtros[1])
	}
	if _, err := consulta.OrderBy("-FechaModificacion", "-Id").All(&solicitudes); err != nil {
		return nil, err
	}
	resultado := make([]models.BorradorInscripcionGrado, 0, limit)
	omitidas := 0
	for _, solicitud := range solicitudes {
		expediente, err := cargarSolicitudRevision(o, solicitud, estadosPermitidos)
		if errors.Is(err, services.ErrBorradorCerrado) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if omitidas < offset {
			omitidas++
			continue
		}
		resultado = append(resultado, *expediente)
		if len(resultado) == limit {
			break
		}
	}
	return resultado, nil
}

func (r *RevisionDocumentalORM) Guardar(id, dependenciaID, estadoSolicitudID int, entrada models.RevisarDocumentacionGrado) (*models.BorradorInscripcionGrado, error) {
	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var solicitud models.SolicitudGrado
	if err := tx.QueryTable(new(models.SolicitudGrado)).ForUpdate().Filter("Id", id).
		Filter("DependenciaOikosId", dependenciaID).Filter("Activo", true).One(&solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, services.ErrSolicitudFueraAlcance
		}
		return nil, err
	}
	resultado, err := cargarSolicitudRevision(tx, solicitud, []int{entrada.EstadoRadicadaId})
	if err != nil {
		return nil, err
	}
	if resultado.Formulario.Id != entrada.FormularioId || resultado.Formulario.FechaRadicacion == nil ||
		len(resultado.Soportes) != len(entrada.Soportes) || len(resultado.Soportes) < 3 || len(resultado.Soportes) > 4 {
		return nil, services.ErrRevisionDocumentalInvalida
	}
	decisiones := make(map[int]models.DecisionSoporteGrado, len(entrada.Soportes))
	for _, decision := range entrada.Soportes {
		decisiones[decision.SoporteId] = decision
	}
	ahora := utils.HoraBogota()
	for _, soporte := range resultado.Soportes {
		decision, existe := decisiones[soporte.Id]
		if !existe || soporte.FormularioSolicitudGradoId != entrada.FormularioId {
			return nil, services.ErrRevisionDocumentalInvalida
		}
		historial := models.HistorialSoporteGrado{
			SoporteGradoId: soporte.Id, TerceroId: entrada.TerceroId,
			EstadoSoporteId: decision.EstadoSoporteId, Observacion: decision.Observacion,
			Activo: true, FechaCreacion: ahora, FechaModificacion: ahora,
		}
		if _, err := tx.Insert(&historial); err != nil {
			return nil, err
		}
	}
	historialSolicitud := models.HistorialSolicitudGrado{
		SolicitudGradoId: id, FormularioSolicitudGradoId: &entrada.FormularioId,
		TerceroId: entrada.TerceroId, EstadoSolicitudId: estadoSolicitudID,
		Justificacion: &entrada.Justificacion, Activo: true,
		FechaCreacion: ahora, FechaModificacion: ahora,
	}
	if _, err := tx.Insert(&historialSolicitud); err != nil {
		return nil, err
	}
	if entrada.Aprobada {
		for _, tipoID := range entrada.TiposPazSalvoId {
			pazSalvo := models.PazSalvoGrado{
				SolicitudGradoId: id, TipoPazSalvoId: tipoID, Activo: true,
				FechaCreacion: ahora, FechaModificacion: ahora,
			}
			if _, err := tx.Insert(&pazSalvo); err != nil {
				return nil, err
			}
			historial := models.HistorialPazSalvoGrado{
				PazSalvoId: pazSalvo.Id, TerceroId: entrada.TerceroId,
				EstadoPazSalvoId: entrada.EstadoPazSalvoPendienteId,
				Justificacion:    "Inicio tras aprobación documental de Secretaría", Activo: true,
				FechaCreacion: ahora, FechaModificacion: ahora,
			}
			if _, err := tx.Insert(&historial); err != nil {
				return nil, err
			}
		}
	}
	resultado, err = cargarSolicitudRevision(tx, solicitud, []int{estadoSolicitudID})
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return resultado, nil
}
