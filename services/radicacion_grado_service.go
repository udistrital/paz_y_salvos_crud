package services

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/utils"
)

// RadicarInscripcionGrado inmoviliza en una transacción el formulario y sus
// soportes, y registra el nuevo estado sin borrar el historial previo.
func RadicarInscripcionGrado(id, terceroID int, entrada models.RadicarInscripcionGrado) (*models.BorradorInscripcionGrado, error) {
	if id <= 0 || terceroID <= 0 || entrada.FormularioId <= 0 || entrada.EstadoBorradorId <= 0 ||
		entrada.EstadoRadicadaId <= 0 || entrada.EstadoRadicadaId == entrada.EstadoBorradorId ||
		entrada.EstadoSoportePendienteId <= 0 || len(entrada.TiposSoporteId) < 3 || len(entrada.TiposSoporteId) > 4 ||
		!contenidoRadicacionValido(entrada.Contenido) {
		return nil, ErrRadicacionInvalida
	}
	tipos := make(map[int]bool, len(entrada.TiposSoporteId))
	for _, tipo := range entrada.TiposSoporteId {
		if tipo <= 0 || tipos[tipo] {
			return nil, ErrRadicacionInvalida
		}
		tipos[tipo] = true
	}
	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	resultado := &models.BorradorInscripcionGrado{}
	if err := tx.Raw(`SELECT * FROM paz_y_salvos.solicitud_grado WHERE id=? AND tercero_id=? AND activo FOR UPDATE`,
		id, terceroID).QueryRow(&resultado.Solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, ErrBorradorNoEncontrado
		}
		return nil, err
	}
	if err := cargarBorrador(tx, resultado, entrada.EstadoBorradorId); err != nil {
		return nil, err
	}
	if resultado.Formulario.Id != entrada.FormularioId {
		return nil, ErrBorradorCerrado
	}
	var conteo struct {
		Total      int `orm:"column(total)"`
		Tipos      int `orm:"column(tipos)"`
		Documentos int `orm:"column(documentos)"`
		Pendientes int `orm:"column(pendientes)"`
	}
	err = tx.Raw(`SELECT COUNT(*) AS total, COUNT(DISTINCT s.tipo_documento_id) AS tipos,
		COUNT(DISTINCT s.documento_id) AS documentos,
		COUNT(*) FILTER (WHERE h.estado_soporte_id=?) AS pendientes
		FROM paz_y_salvos.soporte_solicitud_grado s
		JOIN LATERAL (SELECT estado_soporte_id FROM paz_y_salvos.historial_soporte_solicitud_grado
			WHERE soporte_solicitud_grado_id=s.id AND activo ORDER BY id DESC LIMIT 1) h ON true
		WHERE s.solicitud_grado_id=? AND s.formulario_solicitud_grado_id=? AND s.activo
		`, entrada.EstadoSoportePendienteId, id, entrada.FormularioId).QueryRow(&conteo)
	cantidad := len(entrada.TiposSoporteId)
	if err != nil || conteo.Total != cantidad || conteo.Tipos != cantidad || conteo.Documentos != cantidad ||
		conteo.Pendientes != cantidad || len(resultado.Soportes) != cantidad {
		return nil, ErrRadicacionInvalida
	}
	for _, soporte := range resultado.Soportes {
		if !tipos[soporte.TipoDocumentoId] {
			return nil, ErrRadicacionInvalida
		}
	}
	ahora := utils.HoraBogota()
	if err := tx.Raw(`UPDATE paz_y_salvos.formulario_solicitud_grado
		SET contenido=?::jsonb, fecha_radicacion=?, fecha_modificacion=?
		WHERE id=? AND solicitud_grado_id=? AND fecha_radicacion IS NULL AND activo RETURNING *`,
		string(entrada.Contenido), ahora, ahora, entrada.FormularioId, id).QueryRow(&resultado.Formulario); err != nil {
		return nil, ErrBorradorCerrado
	}
	if err := cargarContenidoFormulario(tx, &resultado.Formulario); err != nil {
		return nil, err
	}
	if err := tx.Raw(`INSERT INTO paz_y_salvos.historial_solicitud_grado
		(solicitud_grado_id, formulario_solicitud_grado_id, tercero_id, estado_solicitud_id, justificacion,
		 fecha_creacion, fecha_modificacion)
		VALUES (?, ?, ?, ?, 'Inscripción radicada por el estudiante', ?, ?) RETURNING *`, id, entrada.FormularioId,
		terceroID, entrada.EstadoRadicadaId, ahora, ahora).QueryRow(&resultado.Historial); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return resultado, nil
}
