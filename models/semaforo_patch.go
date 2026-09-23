package models

import (
	"bytes"
	"encoding/json"
)

type NullableBoolPatch struct {
	Present bool
	Value   *bool
}

func (field *NullableBoolPatch) UnmarshalJSON(data []byte) error {
	field.Present = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		field.Value = nil
		return nil
	}

	var value bool
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil
}

type SemaforoPatch struct {
	Academico               *bool             `json:"Academico"`
	Financiero              *bool             `json:"Financiero"`
	Biblioteca              *bool             `json:"Biblioteca"`
	Laboratorios            *bool             `json:"Laboratorios"`
	Bienestar               *bool             `json:"Bienestar"`
	Urelinter               *bool             `json:"Urelinter"`
	Orc                     NullableBoolPatch `json:"Orc"`
	ObservacionCoordinacion *string           `json:"ObservacionCoordinacion"`
	ObservacionBiblioteca   *string           `json:"ObservacionBiblioteca"`
	ObservacionLaboratorios *string           `json:"ObservacionLaboratorios"`
	ObservacionBienestar    *string           `json:"ObservacionBienestar"`
	ObservacionUrelinter    *string           `json:"ObservacionUrelinter"`
	ObservacionOrc          *string           `json:"ObservacionOrc"`
	ObservacionFinanciera   *string           `json:"ObservacionFinanciera"`
}

func (patch SemaforoPatch) ApplyTo(semaforo *Semaforo) []string {
	columns := make([]string, 0, 14)
	if patch.Academico != nil {
		semaforo.Academico = *patch.Academico
		columns = append(columns, "Academico")
	}
	if patch.Financiero != nil {
		semaforo.Financiero = *patch.Financiero
		columns = append(columns, "Financiero")
	}
	if patch.Biblioteca != nil {
		semaforo.Biblioteca = *patch.Biblioteca
		columns = append(columns, "Biblioteca")
	}
	if patch.Laboratorios != nil {
		semaforo.Laboratorios = *patch.Laboratorios
		columns = append(columns, "Laboratorios")
	}
	if patch.Bienestar != nil {
		semaforo.Bienestar = *patch.Bienestar
		columns = append(columns, "Bienestar")
	}
	if patch.Urelinter != nil {
		semaforo.Urelinter = *patch.Urelinter
		columns = append(columns, "Urelinter")
	}
	if patch.Orc.Present {
		semaforo.Orc = patch.Orc.Value
		columns = append(columns, "Orc")
	}
	if patch.ObservacionCoordinacion != nil {
		semaforo.ObservacionCoordinacion = *patch.ObservacionCoordinacion
		columns = append(columns, "ObservacionCoordinacion")
	}
	if patch.ObservacionBiblioteca != nil {
		semaforo.ObservacionBiblioteca = *patch.ObservacionBiblioteca
		columns = append(columns, "ObservacionBiblioteca")
	}
	if patch.ObservacionLaboratorios != nil {
		semaforo.ObservacionLaboratorios = *patch.ObservacionLaboratorios
		columns = append(columns, "ObservacionLaboratorios")
	}
	if patch.ObservacionBienestar != nil {
		semaforo.ObservacionBienestar = *patch.ObservacionBienestar
		columns = append(columns, "ObservacionBienestar")
	}
	if patch.ObservacionUrelinter != nil {
		semaforo.ObservacionUrelinter = *patch.ObservacionUrelinter
		columns = append(columns, "ObservacionUrelinter")
	}
	if patch.ObservacionOrc != nil {
		semaforo.ObservacionOrc = *patch.ObservacionOrc
		columns = append(columns, "ObservacionOrc")
	}
	if patch.ObservacionFinanciera != nil {
		semaforo.ObservacionFinanciera = *patch.ObservacionFinanciera
		columns = append(columns, "ObservacionFinanciera")
	}
	return columns
}

type DeletedSemaforoResponse struct {
	Id int `json:"Id"`
}
