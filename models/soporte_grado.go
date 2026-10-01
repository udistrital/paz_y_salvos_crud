package models

import "time"

type SoporteGrado struct {
	Id                         int       `json:"Id" orm:"column(id)"`
	SolicitudGradoId           int       `json:"SolicitudGradoId" orm:"column(solicitud_grado_id)"`
	FormularioSolicitudGradoId int       `json:"FormularioSolicitudGradoId" orm:"column(formulario_solicitud_grado_id)"`
	DocumentoId                int       `json:"DocumentoId" orm:"column(documento_id)"`
	TipoDocumentoId            int       `json:"TipoDocumentoId" orm:"column(tipo_documento_id)"`
	SoporteAnteriorId          *int      `json:"SoporteAnteriorId" orm:"column(soporte_anterior_id);null"`
	Activo                     bool      `json:"Activo" orm:"column(activo)"`
	FechaCreacion              time.Time `json:"FechaCreacion" orm:"column(fecha_creacion)"`
	FechaModificacion          time.Time `json:"FechaModificacion" orm:"column(fecha_modificacion)"`
}

// Identidad y parámetros provienen del MID; la versión y referencia previa
// permiten detectar una pantalla obsoleta bajo el bloqueo de la solicitud.
type AsociarSoporteGrado struct {
	FormularioId    int `json:"FormularioId"`
	SoporteActualId int `json:"SoporteActualId"`
	DocumentoId     int `json:"DocumentoId"`
	EstadoSoporteId int `json:"EstadoSoporteId"`
}
