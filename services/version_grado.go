package services

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
)

// El MID envía exclusivamente el ámbito que resolvió desde la sesión.
func solicitudVersionEnAlcance(o orm.QueryExecutor, id, terceroID, dependenciaID int) (models.SolicitudGrado, error) {
	var solicitud models.SolicitudGrado
	if id <= 0 || (terceroID <= 0) == (dependenciaID <= 0) {
		return solicitud, ErrBorradorInvalido
	}
	query := `SELECT * FROM paz_y_salvos.solicitud_grado WHERE id=? AND activo AND `
	valor := terceroID
	if terceroID > 0 {
		query += `tercero_id=?`
	} else {
		query += `dependencia_oikos_id=?`
		valor = dependenciaID
	}
	if err := o.Raw(query, id, valor).QueryRow(&solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return solicitud, ErrBorradorNoEncontrado
		}
		return solicitud, err
	}
	return solicitud, nil
}

func ListarVersionesGrado(id, terceroID, dependenciaID int) ([]models.VersionGradoResumen, error) {
	o := orm.NewOrm()
	if _, err := solicitudVersionEnAlcance(o, id, terceroID, dependenciaID); err != nil {
		return nil, err
	}
	var formularios []models.FormularioSolicitudGrado
	if _, err := o.Raw(`SELECT * FROM paz_y_salvos.formulario_solicitud_grado
		WHERE solicitud_grado_id=? AND activo AND fecha_radicacion IS NOT NULL ORDER BY version`, id).QueryRows(&formularios); err != nil {
		return nil, err
	}
	versiones := make([]models.VersionGradoResumen, 0, len(formularios))
	for _, formulario := range formularios {
		var historial models.HistorialSolicitudGrado
		if err := o.Raw(`SELECT * FROM paz_y_salvos.historial_solicitud_grado
			WHERE solicitud_grado_id=? AND formulario_solicitud_grado_id=? AND activo
			ORDER BY id DESC LIMIT 1`, id, formulario.Id).QueryRow(&historial); err != nil {
			return nil, err
		}
		fecha := formulario.FechaRadicacion.Format("2006-01-02T15:04:05")
		versiones = append(versiones, models.VersionGradoResumen{
			FormularioId: formulario.Id, Version: formulario.Version, FechaRadicacion: fecha,
			EstadoSolicitudId: historial.EstadoSolicitudId, Justificacion: historial.Justificacion,
		})
	}
	return versiones, nil
}

func ConsultarVersionGrado(id, version, terceroID, dependenciaID int) (*models.VersionGradoDetalle, error) {
	if version <= 0 {
		return nil, ErrBorradorInvalido
	}
	o := orm.NewOrm()
	solicitud, err := solicitudVersionEnAlcance(o, id, terceroID, dependenciaID)
	if err != nil {
		return nil, err
	}
	resultado := &models.VersionGradoDetalle{Solicitud: solicitud, Historial: []models.HistorialSolicitudGrado{}, Soportes: []models.SoporteVersionGrado{}}
	if err := o.Raw(`SELECT * FROM paz_y_salvos.formulario_solicitud_grado
		WHERE solicitud_grado_id=? AND version=? AND activo AND fecha_radicacion IS NOT NULL`, id, version).QueryRow(&resultado.Formulario); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, ErrBorradorNoEncontrado
		}
		return nil, err
	}
	if err := cargarContenidoFormulario(o, &resultado.Formulario); err != nil {
		return nil, err
	}
	if _, err := o.Raw(`SELECT * FROM paz_y_salvos.historial_solicitud_grado
		WHERE solicitud_grado_id=? AND formulario_solicitud_grado_id=? AND activo ORDER BY id`, id, resultado.Formulario.Id).
		QueryRows(&resultado.Historial); err != nil {
		return nil, err
	}
	if len(resultado.Historial) == 0 {
		return nil, ErrBorradorCerrado
	}
	var soportes []models.SoporteGrado
	if _, err := o.Raw(`SELECT * FROM paz_y_salvos.soporte_solicitud_grado
		WHERE solicitud_grado_id=? AND formulario_solicitud_grado_id=? AND activo ORDER BY tipo_documento_id`, id, resultado.Formulario.Id).
		QueryRows(&soportes); err != nil {
		return nil, err
	}
	for _, soporte := range soportes {
		item := models.SoporteVersionGrado{Soporte: soporte, Historial: []models.HistorialSoporteGrado{}}
		if _, err := o.Raw(`SELECT * FROM paz_y_salvos.historial_soporte_solicitud_grado
			WHERE soporte_solicitud_grado_id=? AND activo ORDER BY id`, soporte.Id).QueryRows(&item.Historial); err != nil {
			return nil, err
		}
		if len(item.Historial) == 0 {
			return nil, ErrBorradorCerrado
		}
		item.EstadoActual = item.Historial[len(item.Historial)-1]
		resultado.Soportes = append(resultado.Soportes, item)
	}
	return resultado, nil
}
