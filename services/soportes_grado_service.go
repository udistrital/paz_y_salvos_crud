package services

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
)

// AsociarSoporteBorrador nunca sobrescribe el documento anterior ni sus
// actuaciones. Solo puede cambiar la referencia activa de una versión borrador.
func AsociarSoporteBorrador(id, terceroID, estadoBorradorID, tipoID int, entrada models.AsociarSoporteGrado) (*models.BorradorInscripcionGrado, error) {
	if id <= 0 || terceroID <= 0 || estadoBorradorID <= 0 || tipoID <= 0 || entrada.FormularioId <= 0 ||
		entrada.DocumentoId <= 0 || entrada.EstadoSoporteId <= 0 || entrada.SoporteActualId < 0 {
		return nil, ErrBorradorInvalido
	}
	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	borrador := &models.BorradorInscripcionGrado{}
	if err := tx.Raw(`SELECT * FROM paz_y_salvos.solicitud_grado WHERE id=? AND tercero_id=? AND activo FOR UPDATE`,
		id, terceroID).QueryRow(&borrador.Solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, ErrBorradorNoEncontrado
		}
		return nil, err
	}
	if err := cargarBorrador(tx, borrador, estadoBorradorID); err != nil {
		return nil, err
	}
	if borrador.Formulario.Id != entrada.FormularioId {
		return nil, ErrBorradorCerrado
	}
	var anterior *models.SoporteGrado
	for i := range borrador.Soportes {
		s := &borrador.Soportes[i]
		if s.DocumentoId == entrada.DocumentoId && s.TipoDocumentoId != tipoID {
			return nil, fmt.Errorf("%w: cada tipo de soporte requiere una referencia diferente", ErrBorradorInvalido)
		}
		if s.TipoDocumentoId == tipoID {
			if anterior != nil {
				return nil, ErrBorradorCerrado
			}
			anterior = s
		}
	}
	actual := 0
	if anterior != nil {
		actual = anterior.Id
	}
	if actual != entrada.SoporteActualId {
		return nil, ErrBorradorCerrado
	}
	if anterior != nil && anterior.DocumentoId == entrada.DocumentoId {
		return borrador, nil // Reintento de asociación, sin nueva actuación.
	}
	if anterior != nil {
		if _, err := tx.Raw(`UPDATE paz_y_salvos.soporte_solicitud_grado SET activo=false, fecha_modificacion=now() WHERE id=?`, anterior.Id).Exec(); err != nil {
			return nil, err
		}
	}
	var soporte models.SoporteGrado
	if err := tx.Raw(`INSERT INTO paz_y_salvos.soporte_solicitud_grado
		(solicitud_grado_id, formulario_solicitud_grado_id, documento_id, tipo_documento_id, soporte_anterior_id)
		VALUES (?, ?, ?, ?, NULLIF(?, 0)) RETURNING *`, id, entrada.FormularioId, entrada.DocumentoId, tipoID, actual).QueryRow(&soporte); err != nil {
		return nil, err
	}
	if _, err := tx.Raw(`INSERT INTO paz_y_salvos.historial_soporte_solicitud_grado
		(soporte_solicitud_grado_id, tercero_id, estado_soporte_id, observacion)
		VALUES (?, ?, ?, 'Soporte provisional cargado por el estudiante')`, soporte.Id, terceroID, entrada.EstadoSoporteId).Exec(); err != nil {
		return nil, err
	}
	if err := cargarBorrador(tx, borrador, estadoBorradorID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return borrador, nil
}
