package services

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/utils"
)

// PatchSemaforoWithValidation aplica un PATCH con bloqueo de fila para que la
// validación y la actualización operen sobre el mismo estado del registro.
func PatchSemaforoWithValidation(id int, patch models.SemaforoPatch) error {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return errors.New("Error al iniciar transacción")
	}

	currentRecord := &models.Semaforo{}
	if err = tx.Raw("SELECT * FROM semaforo WHERE id = ? FOR UPDATE", id).QueryRow(currentRecord); err != nil {
		tx.Rollback()
		return errors.New("Registro no encontrado")
	}

	if patch.Orc.Present {
		if err = validateOrcUpdate(currentRecord, patch.Orc.Value); err != nil {
			tx.Rollback()
			return err
		}
	}

	if currentRecord.Orc != nil {
		if err = validateFieldChangesWhenOrcActive(currentRecord, patch); err != nil {
			tx.Rollback()
			return err
		}
	}

	columns := patch.ApplyTo(currentRecord)
	if len(columns) == 0 {
		tx.Rollback()
		return errors.New("No se enviaron campos para actualizar")
	}
	currentRecord.FechaModificacion = utils.HoraBogota()
	columns = append(columns, "FechaModificacion")
	if _, err = tx.Update(currentRecord, columns...); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func validateOrcUpdate(currentRecord *models.Semaforo, newOrc *bool) error {
	if newOrc != nil && *newOrc {
		if !currentRecord.Academico || !currentRecord.Financiero ||
			!currentRecord.Biblioteca || !currentRecord.Laboratorios ||
			!currentRecord.Bienestar || !currentRecord.Urelinter {
			return errors.New("No se puede establecer ORC en true: existen dependencias en false")
		}
	}

	if newOrc != nil && !*newOrc {
		if currentRecord.Academico && currentRecord.Financiero &&
			currentRecord.Biblioteca && currentRecord.Laboratorios &&
			currentRecord.Bienestar && currentRecord.Urelinter {
			return errors.New("No se puede establecer ORC en false: todas las dependencias están en true")
		}
	}
	return nil
}

func validateFieldChangesWhenOrcActive(current *models.Semaforo, patch models.SemaforoPatch) error {
	changesState :=
		(patch.Academico != nil && *patch.Academico != current.Academico) ||
			(patch.Financiero != nil && *patch.Financiero != current.Financiero) ||
			(patch.Biblioteca != nil && *patch.Biblioteca != current.Biblioteca) ||
			(patch.Laboratorios != nil && *patch.Laboratorios != current.Laboratorios) ||
			(patch.Bienestar != nil && *patch.Bienestar != current.Bienestar) ||
			(patch.Urelinter != nil && *patch.Urelinter != current.Urelinter)

	if changesState {
		return errors.New("No se pueden modificar campos cuando ORC no es neutro. Solo se permite actualizar el campo ORC y observaciones")
	}
	return nil
}
