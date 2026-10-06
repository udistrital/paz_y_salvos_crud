package models

import (
	"encoding/json"
	"time"
)

type CrearBorradorInscripcionGrado struct {
	TerceroId                     int             `json:"TerceroId"`
	CodigoEstudiante              string          `json:"CodigoEstudiante"`
	PeriodoId                     int             `json:"PeriodoId"`
	ProgramaAcademicoId           int             `json:"ProgramaAcademicoId"`
	DependenciaOikosId            int             `json:"DependenciaOikosId"`
	CalendarioEventoInscripcionId int             `json:"CalendarioEventoInscripcionId"`
	CalendarioEventoAprobacionId  int             `json:"CalendarioEventoAprobacionId"`
	EstadoBorradorId              int             `json:"EstadoBorradorId"`
	Contenido                     json.RawMessage `json:"Contenido"`
}

// ActualizarBorradorInscripcionGrado solo acepta el contenido editable; los IDs
// de identidad y estado provienen del MID y se verifican en la operación.
type ActualizarBorradorInscripcionGrado struct {
	Contenido json.RawMessage `json:"Contenido"`
}

type RadicarInscripcionGrado struct {
	FormularioId             int             `json:"FormularioId"`
	Contenido                json.RawMessage `json:"Contenido"`
	EstadoBorradorId         int             `json:"EstadoBorradorId"`
	EstadoRadicadaId         int             `json:"EstadoRadicadaId"`
	EstadoSoportePendienteId int             `json:"EstadoSoportePendienteId"`
	TiposSoporteId           []int           `json:"TiposSoporteId"`
}

type SubsanarInscripcionGrado struct {
	FormularioId             int `json:"FormularioId"`
	EstadoObservadaId        int `json:"EstadoObservadaId"`
	EstadoBorradorId         int `json:"EstadoBorradorId"`
	EstadoSoportePendienteId int `json:"EstadoSoportePendienteId"`
}

// RevisarDocumentacionGrado recibe exclusivamente IDs ya resueltos por el MID.
// El CRUD vuelve a comprobar estado, versión, soportes y atomicidad.
type RevisarDocumentacionGrado struct {
	FormularioId                  int                    `json:"FormularioId"`
	TerceroId                     int                    `json:"TerceroId"`
	Aprobada                      bool                   `json:"Aprobada"`
	Justificacion                 string                 `json:"Justificacion"`
	EstadoRadicadaId              int                    `json:"EstadoRadicadaId"`
	EstadoObservadaId             int                    `json:"EstadoObservadaId"`
	EstadoDocumentacionAprobadaId int                    `json:"EstadoDocumentacionAprobadaId"`
	EstadoSoporteObservadoId      int                    `json:"EstadoSoporteObservadoId"`
	EstadoSoporteAprobadoId       int                    `json:"EstadoSoporteAprobadoId"`
	EstadoPazSalvoPendienteId     int                    `json:"EstadoPazSalvoPendienteId"`
	TiposPazSalvoId               []int                  `json:"TiposPazSalvoId"`
	Soportes                      []DecisionSoporteGrado `json:"Soportes"`
}

type SolicitudGrado struct {
	Id                            int       `json:"Id" orm:"column(id);pk;auto"`
	TerceroId                     int       `json:"TerceroId" orm:"column(tercero_id)"`
	CodigoEstudiante              string    `json:"CodigoEstudiante" orm:"column(codigo_estudiante)"`
	PeriodoId                     int       `json:"PeriodoId" orm:"column(periodo_id)"`
	ProgramaAcademicoId           int       `json:"ProgramaAcademicoId" orm:"column(programa_academico_id)"`
	DependenciaOikosId            int       `json:"DependenciaOikosId" orm:"column(dependencia_oikos_id)"`
	CalendarioEventoInscripcionId int       `json:"CalendarioEventoInscripcionId" orm:"column(calendario_evento_inscripcion_id)"`
	CalendarioEventoAprobacionId  int       `json:"CalendarioEventoAprobacionId" orm:"column(calendario_evento_aprobacion_id)"`
	Activo                        bool      `json:"Activo" orm:"column(activo)"`
	FechaCreacion                 time.Time `json:"FechaCreacion" orm:"column(fecha_creacion)"`
	FechaModificacion             time.Time `json:"FechaModificacion" orm:"column(fecha_modificacion)"`
}

type FormularioSolicitudGrado struct {
	Id                int             `json:"Id" orm:"column(id);pk;auto"`
	SolicitudGradoId  int             `json:"SolicitudGradoId" orm:"column(solicitud_grado_id)"`
	Version           int16           `json:"Version" orm:"column(version)"`
	Contenido         json.RawMessage `json:"Contenido" orm:"column(contenido);type(json)"`
	Activo            bool            `json:"Activo" orm:"column(activo)"`
	FechaCreacion     time.Time       `json:"FechaCreacion" orm:"column(fecha_creacion)"`
	FechaModificacion time.Time       `json:"FechaModificacion" orm:"column(fecha_modificacion)"`
	FechaRadicacion   *time.Time      `json:"FechaRadicacion" orm:"column(fecha_radicacion);null"`
}

type HistorialSolicitudGrado struct {
	Id                         int       `json:"Id" orm:"column(id);pk;auto"`
	SolicitudGradoId           int       `json:"SolicitudGradoId" orm:"column(solicitud_grado_id)"`
	FormularioSolicitudGradoId *int      `json:"FormularioSolicitudGradoId" orm:"column(formulario_solicitud_grado_id);null"`
	TerceroId                  int       `json:"TerceroId" orm:"column(tercero_id)"`
	EstadoSolicitudId          int       `json:"EstadoSolicitudId" orm:"column(estado_solicitud_id)"`
	Justificacion              *string   `json:"Justificacion" orm:"column(justificacion);null"`
	Activo                     bool      `json:"Activo" orm:"column(activo)"`
	FechaCreacion              time.Time `json:"FechaCreacion" orm:"column(fecha_creacion)"`
	FechaModificacion          time.Time `json:"FechaModificacion" orm:"column(fecha_modificacion)"`
}

type BorradorInscripcionGrado struct {
	Solicitud       SolicitudGrado           `json:"Solicitud"`
	Formulario      FormularioSolicitudGrado `json:"Formulario"`
	Historial       HistorialSolicitudGrado  `json:"Historial"`
	Soportes        []SoporteGrado           `json:"Soportes"`
	EstadosSoportes []HistorialSoporteGrado  `json:"EstadosSoportes"`
}

func (SolicitudGrado) TableName() string           { return "solicitud_grado" }
func (FormularioSolicitudGrado) TableName() string { return "formulario_solicitud_grado" }
func (HistorialSolicitudGrado) TableName() string  { return "historial_solicitud_grado" }
