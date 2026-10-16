package repositories

import (
	"errors"
	"sort"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/services"
	"github.com/udistrital/paz_y_salvos_crud/utils"
)

type PazSalvoGradoORM struct{}

func NuevoPazSalvoGradoORM() *PazSalvoGradoORM { return &PazSalvoGradoORM{} }

func contieneID(ids []int, esperado int) bool {
	for _, id := range ids {
		if id == esperado {
			return true
		}
	}
	return false
}

func estadoSolicitudActual(o orm.QueryExecutor, solicitudID int) (models.HistorialSolicitudGrado, error) {
	var historial models.HistorialSolicitudGrado
	err := o.QueryTable(new(models.HistorialSolicitudGrado)).Filter("SolicitudGradoId", solicitudID).
		Filter("Activo", true).OrderBy("-FechaCreacion", "-Id").One(&historial)
	return historial, err
}

func cargarPazSalvosGrado(o orm.QueryExecutor, solicitudID, estadoDocumentacionAprobadaID int, tiposEsperados []int) (*models.PazSalvosSolicitudGrado, error) {
	var solicitud models.SolicitudGrado
	if err := o.QueryTable(new(models.SolicitudGrado)).Filter("Id", solicitudID).Filter("Activo", true).One(&solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, services.ErrPazSalvoGradoNoEncontrado
		}
		return nil, err
	}
	historialSolicitud, err := estadoSolicitudActual(o, solicitudID)
	if err != nil || historialSolicitud.EstadoSolicitudId != estadoDocumentacionAprobadaID {
		if err != nil && !errors.Is(err, orm.ErrNoRows) {
			return nil, err
		}
		return nil, services.ErrPazSalvoGradoConflicto
	}
	var filas []models.PazSalvoGrado
	if _, err := o.QueryTable(new(models.PazSalvoGrado)).Filter("SolicitudGradoId", solicitudID).
		Filter("Activo", true).All(&filas); err != nil {
		return nil, err
	}
	if len(filas) != len(tiposEsperados) {
		return nil, services.ErrPazSalvoGradoConflicto
	}
	porTipo := make(map[int]models.PazSalvoGrado, len(filas))
	for _, fila := range filas {
		if !contieneID(tiposEsperados, fila.TipoPazSalvoId) || porTipo[fila.TipoPazSalvoId].Id != 0 {
			return nil, services.ErrPazSalvoGradoConflicto
		}
		porTipo[fila.TipoPazSalvoId] = fila
	}
	resultado := &models.PazSalvosSolicitudGrado{Solicitud: solicitud, Checks: make([]models.CheckPazSalvoGrado, 0, len(filas))}
	for _, tipo := range tiposEsperados {
		fila := porTipo[tipo]
		var historial []models.HistorialPazSalvoGrado
		if _, err := o.QueryTable(new(models.HistorialPazSalvoGrado)).Filter("PazSalvoId", fila.Id).
			Filter("Activo", true).OrderBy("FechaCreacion", "Id").All(&historial); err != nil {
			return nil, err
		}
		if len(historial) == 0 {
			return nil, services.ErrPazSalvoGradoConflicto
		}
		resultado.Checks = append(resultado.Checks, models.CheckPazSalvoGrado{
			PazSalvo: fila, EstadoActual: historial[len(historial)-1], Historial: historial,
		})
	}
	return resultado, nil
}

func (r *PazSalvoGradoORM) Consultar(solicitudID, estadoDocumentacionAprobadaID int, tiposEsperados []int) (*models.PazSalvosSolicitudGrado, error) {
	return cargarPazSalvosGrado(orm.NewOrm(), solicitudID, estadoDocumentacionAprobadaID, tiposEsperados)
}

func (r *PazSalvoGradoORM) Listar(estadoDocumentacionAprobadaID int, tiposEsperados, dependencias []int, limit, offset, periodoID, programaID, terceroID int, codigo string) (*models.PaginaPazSalvosGrado, error) {
	o := orm.NewOrm()
	consulta := o.QueryTable(new(models.SolicitudGrado)).Filter("Activo", true)
	if len(dependencias) > 0 {
		consulta = consulta.Filter("DependenciaOikosId__in", dependencias)
	}
	if periodoID > 0 {
		consulta = consulta.Filter("PeriodoId", periodoID)
	}
	if programaID > 0 {
		consulta = consulta.Filter("ProgramaAcademicoId", programaID)
	}
	if terceroID > 0 {
		consulta = consulta.Filter("TerceroId", terceroID)
	}
	if codigo != "" {
		consulta = consulta.Filter("CodigoEstudiante", codigo)
	}
	var solicitudes []models.SolicitudGrado
	if _, err := consulta.OrderBy("-FechaModificacion", "-Id").All(&solicitudes); err != nil {
		return nil, err
	}
	resultado := &models.PaginaPazSalvosGrado{Solicitudes: []models.PazSalvosSolicitudGrado{}}
	for _, solicitud := range solicitudes {
		historial, err := estadoSolicitudActual(o, solicitud.Id)
		if errors.Is(err, orm.ErrNoRows) || (err == nil && historial.EstadoSolicitudId != estadoDocumentacionAprobadaID) {
			continue
		}
		if err != nil {
			return nil, err
		}
		expediente, err := cargarPazSalvosGrado(o, solicitud.Id, estadoDocumentacionAprobadaID, tiposEsperados)
		if err != nil {
			return nil, err
		}
		if resultado.Total >= offset && len(resultado.Solicitudes) < limit {
			resultado.Solicitudes = append(resultado.Solicitudes, *expediente)
		}
		resultado.Total++
	}
	return resultado, nil
}

func (r *PazSalvoGradoORM) ListarAprobadosSecretaria(estadoDocumentacionAprobadaID, tipoSecretariaID, estadoAprobadoID int,
	dependencias []int, periodoID, limit, offset int) (*models.PaginaAprobadosSecretariaGrado, error) {
	o := orm.NewOrm()
	// El filtro usa el último estado activo del check de Secretaría y de la
	// solicitud, de modo que Total y cada página incluyan solo aprobaciones vigentes.
	filtro := ` FROM paz_y_salvos.solicitud_grado s
		JOIN paz_y_salvos.paz_salvo ps ON ps.solicitud_grado_id=s.id AND ps.activo=true AND ps.tipo_paz_salvo_id=?
		JOIN LATERAL (
			SELECT h.estado_paz_salvo_id, h.fecha_creacion
			FROM paz_y_salvos.historial_paz_salvo h
			WHERE h.paz_salvo_id=ps.id AND h.activo=true
			ORDER BY h.fecha_creacion DESC, h.id DESC LIMIT 1
		) hp ON hp.estado_paz_salvo_id=?
		JOIN LATERAL (
			SELECT h.estado_solicitud_id FROM paz_y_salvos.historial_solicitud_grado h
			WHERE h.solicitud_grado_id=s.id AND h.activo=true
			ORDER BY h.fecha_creacion DESC, h.id DESC LIMIT 1
		) hs ON hs.estado_solicitud_id=?
		WHERE s.activo=true AND s.periodo_id=?`
	args := []interface{}{tipoSecretariaID, estadoAprobadoID, estadoDocumentacionAprobadaID, periodoID}
	if len(dependencias) > 0 {
		marcadores := make([]string, len(dependencias))
		for i, dependencia := range dependencias {
			marcadores[i] = "?"
			args = append(args, dependencia)
		}
		filtro += " AND s.dependencia_oikos_id IN (" + strings.Join(marcadores, ",") + ")"
	}
	resultado := &models.PaginaAprobadosSecretariaGrado{Solicitudes: []models.SolicitudAprobadaSecretariaGrado{}}
	if err := o.Raw("SELECT COUNT(*)"+filtro, args...).QueryRow(&resultado.Total); err != nil {
		return nil, err
	}
	if offset >= resultado.Total {
		return resultado, nil
	}
	consulta := `SELECT s.id AS solicitud_id, s.tercero_id, s.codigo_estudiante, s.periodo_id,
		s.programa_academico_id, s.dependencia_oikos_id, hp.fecha_creacion AS fecha_aprobacion` + filtro +
		` ORDER BY s.fecha_modificacion DESC, s.id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	if _, err := o.Raw(consulta, args...).QueryRows(&resultado.Solicitudes); err != nil {
		return nil, err
	}
	return resultado, nil
}

func (r *PazSalvoGradoORM) Decidir(solicitudID int, entrada models.DecidirPazSalvoGrado) (*models.PazSalvosSolicitudGrado, error) {
	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var solicitud models.SolicitudGrado
	if err := tx.QueryTable(new(models.SolicitudGrado)).ForUpdate().Filter("Id", solicitudID).Filter("Activo", true).One(&solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, services.ErrPazSalvoGradoNoEncontrado
		}
		return nil, err
	}
	tiposEsperados := append(append([]int{}, entrada.TiposPreviosId...), entrada.TipoSecretariaId)
	sort.Ints(tiposEsperados)
	resultado, err := cargarPazSalvosGrado(tx, solicitudID, entrada.EstadoDocumentacionAprobadaId, tiposEsperados)
	if err != nil {
		return nil, err
	}
	porTipo := make(map[int]models.CheckPazSalvoGrado, len(resultado.Checks))
	for _, check := range resultado.Checks {
		porTipo[check.PazSalvo.TipoPazSalvoId] = check
	}
	check, existe := porTipo[entrada.TipoPazSalvoId]
	if !existe || (check.EstadoActual.EstadoPazSalvoId != entrada.EstadoPendienteId &&
		check.EstadoActual.EstadoPazSalvoId != entrada.EstadoAprobadoId &&
		check.EstadoActual.EstadoPazSalvoId != entrada.EstadoDesaprobadoId) ||
		check.EstadoActual.EstadoPazSalvoId == entrada.EstadoPazSalvoId {
		return nil, services.ErrPazSalvoGradoConflicto
	}
	ahora := utils.HoraBogota()
	historial := models.HistorialPazSalvoGrado{
		PazSalvoId: check.PazSalvo.Id, TerceroId: entrada.TerceroId,
		EstadoPazSalvoId: entrada.EstadoPazSalvoId, Justificacion: entrada.Justificacion,
		Activo: true, FechaCreacion: ahora, FechaModificacion: ahora,
	}
	if _, err := tx.Insert(&historial); err != nil {
		return nil, err
	}
	resultado, err = cargarPazSalvosGrado(tx, solicitudID, entrada.EstadoDocumentacionAprobadaId, tiposEsperados)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return resultado, nil
}
