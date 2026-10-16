package models

// VersionGradoResumen contiene únicamente metadatos de una entrega radicada.
type VersionGradoResumen struct {
	FormularioId      int     `json:"FormularioId"`
	Version           int16   `json:"Version"`
	FechaRadicacion   string  `json:"FechaRadicacion"`
	EstadoSolicitudId int     `json:"EstadoSolicitudId"`
	Justificacion     *string `json:"Justificacion"`
}

type SoporteVersionGrado struct {
	Soporte      SoporteGrado            `json:"Soporte"`
	EstadoActual HistorialSoporteGrado   `json:"EstadoActual"`
	Historial    []HistorialSoporteGrado `json:"Historial"`
}

type VersionGradoDetalle struct {
	Solicitud  SolicitudGrado            `json:"Solicitud"`
	Formulario FormularioSolicitudGrado  `json:"Formulario"`
	Historial  []HistorialSolicitudGrado `json:"Historial"`
	Soportes   []SoporteVersionGrado     `json:"Soportes"`
}
