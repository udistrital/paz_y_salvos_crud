package repositories

import (
	"errors"
	"sort"

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
	final := porTipo[entrada.TipoSecretariaId]
	if entrada.TipoPazSalvoId == entrada.TipoSecretariaId {
		if entrada.EstadoPazSalvoId == entrada.EstadoDesaprobadoId {
			if final.EstadoActual.EstadoPazSalvoId != entrada.EstadoAprobadoId {
				return nil, services.ErrPazSalvoGradoConflicto
			}
		} else {
			for _, tipo := range entrada.TiposPreviosId {
				previo, ok := porTipo[tipo]
				if !ok || previo.EstadoActual.EstadoPazSalvoId != entrada.EstadoAprobadoId {
					return nil, services.ErrPazSalvoGradoConflicto
				}
			}
		}
	} else if final.EstadoActual.EstadoPazSalvoId == entrada.EstadoAprobadoId {
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
