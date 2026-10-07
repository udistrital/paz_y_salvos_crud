package models

import "time"

type SoporteGrado struct {
	Id                         int       `json:"Id" orm:"column(id);pk;auto"`
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

type DecisionSoporteGrado struct {
	SoporteId       int    `json:"SoporteId"`
	EstadoSoporteId int    `json:"EstadoSoporteId"`
	Observacion     string `json:"Observacion"`
}

type HistorialSoporteGrado struct {
	Id                int       `orm:"column(id);pk;auto"`
	SoporteGradoId    int       `orm:"column(soporte_solicitud_grado_id)"`
	TerceroId         int       `orm:"column(tercero_id)"`
	EstadoSoporteId   int       `orm:"column(estado_soporte_id)"`
	Observacion       string    `orm:"column(observacion);null"`
	Activo            bool      `orm:"column(activo)"`
	FechaCreacion     time.Time `orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}

type PazSalvoGrado struct {
	Id                int       `json:"Id" orm:"column(id);pk;auto"`
	SolicitudGradoId  int       `json:"SolicitudGradoId" orm:"column(solicitud_grado_id)"`
	TipoPazSalvoId    int       `json:"TipoPazSalvoId" orm:"column(tipo_paz_salvo_id)"`
	Activo            bool      `json:"Activo" orm:"column(activo)"`
	FechaCreacion     time.Time `json:"FechaCreacion" orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `json:"FechaModificacion" orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}

type HistorialPazSalvoGrado struct {
	Id                int       `json:"Id" orm:"column(id);pk;auto"`
	PazSalvoId        int       `json:"PazSalvoId" orm:"column(paz_salvo_id)"`
	TerceroId         int       `json:"TerceroId" orm:"column(tercero_id)"`
	EstadoPazSalvoId  int       `json:"EstadoPazSalvoId" orm:"column(estado_paz_salvo_id)"`
	Justificacion     string    `json:"Justificacion" orm:"column(justificacion);null"`
	Activo            bool      `json:"Activo" orm:"column(activo)"`
	FechaCreacion     time.Time `json:"FechaCreacion" orm:"column(fecha_creacion);type(timestamp without time zone)"`
	FechaModificacion time.Time `json:"FechaModificacion" orm:"column(fecha_modificacion);type(timestamp without time zone)"`
}

type CheckPazSalvoGrado struct {
	PazSalvo     PazSalvoGrado            `json:"PazSalvo"`
	EstadoActual HistorialPazSalvoGrado   `json:"EstadoActual"`
	Historial    []HistorialPazSalvoGrado `json:"Historial"`
}

type PazSalvosSolicitudGrado struct {
	Solicitud SolicitudGrado       `json:"Solicitud"`
	Checks    []CheckPazSalvoGrado `json:"Checks"`
}

type PaginaPazSalvosGrado struct {
	Solicitudes []PazSalvosSolicitudGrado `json:"Solicitudes"`
	Total       int                       `json:"Total"`
}

// DecidirPazSalvoGrado es un contrato interno: el MID resuelve los códigos de
// Parámetros y el CRUD comprueba nuevamente el conjunto y la transición.
type DecidirPazSalvoGrado struct {
	TerceroId                     int    `json:"TerceroId"`
	TipoPazSalvoId                int    `json:"TipoPazSalvoId"`
	EstadoPazSalvoId              int    `json:"EstadoPazSalvoId"`
	Justificacion                 string `json:"Justificacion"`
	EstadoDocumentacionAprobadaId int    `json:"EstadoDocumentacionAprobadaId"`
	EstadoPendienteId             int    `json:"EstadoPendienteId"`
	EstadoAprobadoId              int    `json:"EstadoAprobadoId"`
	EstadoDesaprobadoId           int    `json:"EstadoDesaprobadoId"`
	TipoSecretariaId              int    `json:"TipoSecretariaId"`
	TiposPreviosId                []int  `json:"TiposPreviosId"`
}

func (SoporteGrado) TableName() string           { return "soporte_solicitud_grado" }
func (HistorialSoporteGrado) TableName() string  { return "historial_soporte_solicitud_grado" }
func (PazSalvoGrado) TableName() string          { return "paz_salvo" }
func (HistorialPazSalvoGrado) TableName() string { return "historial_paz_salvo" }
