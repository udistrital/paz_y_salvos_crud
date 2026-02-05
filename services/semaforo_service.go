package services

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/udistrital/paz_y_salvos_crud/models"
)

// PatchSemaforoWithValidation aplica un PATCH con validaciones de reglas de negocio
// y protección contra condiciones de carrera mediante transacciones con SELECT FOR UPDATE
func PatchSemaforoWithValidation(id int, params map[string]interface{}) error {
	o := orm.NewOrm()

	// Iniciar transacción para garantizar atomicidad
	err := o.Begin()
	if err != nil {
		return errors.New("Error al iniciar transacción")
	}

	// Obtener el estado actual del registro con bloqueo exclusivo (SELECT FOR UPDATE)
	// Esto previene condiciones de carrera bloqueando la fila hasta el commit
	currentRecord := &models.Semaforo{}
	err = o.Raw("SELECT * FROM semaforo WHERE id = ? FOR UPDATE", id).QueryRow(currentRecord)
	if err != nil {
		o.Rollback()
		return errors.New("Registro no encontrado")
	}

	// Validación 1: Si se está intentando actualizar ORC
	if orcValue, orcExists := params["Orc"]; orcExists {
		err = validateOrcUpdate(currentRecord, orcValue)
		if err != nil {
			o.Rollback()
			return err
		}
	}

	// Validación 2: Si ORC está activo (no null), solo se puede actualizar el propio campo ORC
	if currentRecord.Orc != nil {
		err = validateFieldChangesWhenOrcActive(currentRecord, params)
		if err != nil {
			o.Rollback()
			return err
		}
	}

	// Si pasó todas las validaciones, aplicar el patch dentro de la transacción
	qs := o.QueryTable(new(models.Semaforo))
	_, err = qs.Filter("id", id).Update(params)
	if err != nil {
		o.Rollback()
		return err
	}

	// Commit de la transacción - libera el bloqueo
	return o.Commit()
}

// validateOrcUpdate valida las reglas de negocio al actualizar ORC
func validateOrcUpdate(currentRecord *models.Semaforo, orcValue interface{}) error {
	// Convertir el valor de ORC a *bool
	var newOrc *bool
	if orcValue == nil {
		newOrc = nil
	} else {
		switch v := orcValue.(type) {
		case bool:
			newOrc = &v
		case *bool:
			newOrc = v
		default:
			return errors.New("Valor de ORC inválido")
		}
	}

	// Si ORC se está estableciendo en true, todas las dependencias deben estar en true
	if newOrc != nil && *newOrc {
		if !currentRecord.Academico || !currentRecord.Financiero ||
			!currentRecord.Biblioteca || !currentRecord.Laboratorios ||
			!currentRecord.Bienestar || !currentRecord.Urelinter {
			return errors.New("No se puede establecer ORC en true: existen dependencias en false")
		}
	}

	// Si ORC se está estableciendo en false, al menos una dependencia debe estar en false
	if newOrc != nil && !*newOrc {
		if currentRecord.Academico && currentRecord.Financiero &&
			currentRecord.Biblioteca && currentRecord.Laboratorios &&
			currentRecord.Bienestar && currentRecord.Urelinter {
			return errors.New("No se puede establecer ORC en false: todas las dependencias están en true")
		}
	}

	return nil
}

// validateFieldChangesWhenOrcActive verifica que cuando ORC está activo solo se modifique ORC
func validateFieldChangesWhenOrcActive(currentRecord *models.Semaforo, params map[string]interface{}) error {
	// Verificar qué campos están realmente cambiando comparando con el registro actual
	for field, newValue := range params {
		var currentValue interface{}

		switch field {
		case "Orc":
			currentValue = currentRecord.Orc
		case "Academico":
			currentValue = currentRecord.Academico
		case "Financiero":
			currentValue = currentRecord.Financiero
		case "Biblioteca":
			currentValue = currentRecord.Biblioteca
		case "Laboratorios":
			currentValue = currentRecord.Laboratorios
		case "Bienestar":
			currentValue = currentRecord.Bienestar
		case "Urelinter":
			currentValue = currentRecord.Urelinter
		case "ObservacionCoordinacion":
			currentValue = currentRecord.ObservacionCoordinacion
		case "ObservacionBiblioteca":
			currentValue = currentRecord.ObservacionBiblioteca
		case "ObservacionLaboratorios":
			currentValue = currentRecord.ObservacionLaboratorios
		case "ObservacionBienestar":
			currentValue = currentRecord.ObservacionBienestar
		case "ObservacionUrelinter":
			currentValue = currentRecord.ObservacionUrelinter
		case "ObservacionOrc":
			currentValue = currentRecord.ObservacionOrc
		default:
			continue
		}

		// Comparar valores para detectar si hay un cambio real
		isChanging := isFieldChanging(field, currentValue, newValue, currentRecord)

		// Si el campo está cambiando y no es ORC, bloquear
		if isChanging && field != "Orc" {
			return errors.New("No se pueden modificar campos cuando ORC no es neutro. Solo se permite actualizar el campo ORC")
		}
	}

	return nil
}

// isFieldChanging compara el valor actual con el nuevo para detectar cambios reales
func isFieldChanging(field string, currentValue, newValue interface{}, currentRecord *models.Semaforo) bool {
	if field == "Orc" {
		// Comparación especial para punteros
		if newValue == nil && currentRecord.Orc != nil {
			return true
		} else if newValue != nil && currentRecord.Orc == nil {
			return true
		} else if newValue != nil && currentRecord.Orc != nil {
			newBool, ok := newValue.(bool)
			if ok && newBool != *currentRecord.Orc {
				return true
			}
		}
		return false
	}

	// Para otros campos, comparación directa
	return currentValue != newValue
}
