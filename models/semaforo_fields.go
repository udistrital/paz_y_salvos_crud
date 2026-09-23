package models

import (
	"fmt"
	"time"
)

// SemaforoFields representa una fila cuando el cliente solicita una selección
// parcial mediante el parámetro fields.
type SemaforoFields struct {
	Id                      *int       `json:"Id,omitempty"`
	CodigoEstudiante        *float64   `json:"CodigoEstudiante,omitempty"`
	IdFacultadOikos         *int16     `json:"IdFacultadOikos,omitempty"`
	IdProyectoOikos         *int16     `json:"IdProyectoOikos,omitempty"`
	IdFacultadGedep         *int16     `json:"IdFacultadGedep,omitempty"`
	IdProyectoAccra         *int16     `json:"IdProyectoAccra,omitempty"`
	AnioInsGrado            *float64   `json:"AnioInsGrado,omitempty"`
	PerInsGrado             *float64   `json:"PerInsGrado,omitempty"`
	Academico               *bool      `json:"Academico,omitempty"`
	Financiero              *bool      `json:"Financiero,omitempty"`
	Biblioteca              *bool      `json:"Biblioteca,omitempty"`
	Laboratorios            *bool      `json:"Laboratorios,omitempty"`
	Bienestar               *bool      `json:"Bienestar,omitempty"`
	Urelinter               *bool      `json:"Urelinter,omitempty"`
	Orc                     **bool     `json:"Orc,omitempty"`
	ObservacionCoordinacion *string    `json:"ObservacionCoordinacion,omitempty"`
	ObservacionBiblioteca   *string    `json:"ObservacionBiblioteca,omitempty"`
	ObservacionLaboratorios *string    `json:"ObservacionLaboratorios,omitempty"`
	ObservacionBienestar    *string    `json:"ObservacionBienestar,omitempty"`
	ObservacionUrelinter    *string    `json:"ObservacionUrelinter,omitempty"`
	ObservacionOrc          *string    `json:"ObservacionOrc,omitempty"`
	ObservacionFinanciera   *string    `json:"ObservacionFinanciera,omitempty"`
	Activo                  *bool      `json:"Activo,omitempty"`
	FechaCreacion           *time.Time `json:"FechaCreacion,omitempty"`
	FechaModificacion       *time.Time `json:"FechaModificacion,omitempty"`
}

func selectSemaforoFields(semaforo Semaforo, fields []string) (SemaforoFields, error) {
	var selected SemaforoFields
	for _, field := range fields {
		switch field {
		case "Id":
			selected.Id = &semaforo.Id
		case "CodigoEstudiante":
			selected.CodigoEstudiante = &semaforo.CodigoEstudiante
		case "IdFacultadOikos":
			selected.IdFacultadOikos = &semaforo.IdFacultadOikos
		case "IdProyectoOikos":
			selected.IdProyectoOikos = &semaforo.IdProyectoOikos
		case "IdFacultadGedep":
			selected.IdFacultadGedep = &semaforo.IdFacultadGedep
		case "IdProyectoAccra":
			selected.IdProyectoAccra = &semaforo.IdProyectoAccra
		case "AnioInsGrado":
			selected.AnioInsGrado = &semaforo.AnioInsGrado
		case "PerInsGrado":
			selected.PerInsGrado = &semaforo.PerInsGrado
		case "Academico":
			selected.Academico = &semaforo.Academico
		case "Financiero":
			selected.Financiero = &semaforo.Financiero
		case "Biblioteca":
			selected.Biblioteca = &semaforo.Biblioteca
		case "Laboratorios":
			selected.Laboratorios = &semaforo.Laboratorios
		case "Bienestar":
			selected.Bienestar = &semaforo.Bienestar
		case "Urelinter":
			selected.Urelinter = &semaforo.Urelinter
		case "Orc":
			selected.Orc = &semaforo.Orc
		case "ObservacionCoordinacion":
			selected.ObservacionCoordinacion = &semaforo.ObservacionCoordinacion
		case "ObservacionBiblioteca":
			selected.ObservacionBiblioteca = &semaforo.ObservacionBiblioteca
		case "ObservacionLaboratorios":
			selected.ObservacionLaboratorios = &semaforo.ObservacionLaboratorios
		case "ObservacionBienestar":
			selected.ObservacionBienestar = &semaforo.ObservacionBienestar
		case "ObservacionUrelinter":
			selected.ObservacionUrelinter = &semaforo.ObservacionUrelinter
		case "ObservacionOrc":
			selected.ObservacionOrc = &semaforo.ObservacionOrc
		case "ObservacionFinanciera":
			selected.ObservacionFinanciera = &semaforo.ObservacionFinanciera
		case "Activo":
			selected.Activo = &semaforo.Activo
		case "FechaCreacion":
			selected.FechaCreacion = &semaforo.FechaCreacion
		case "FechaModificacion":
			selected.FechaModificacion = &semaforo.FechaModificacion
		default:
			return SemaforoFields{}, fmt.Errorf("campo de selección no válido: %s", field)
		}
	}
	return selected, nil
}
