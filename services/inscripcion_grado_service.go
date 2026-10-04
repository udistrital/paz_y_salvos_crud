package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/lib/pq"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/utils"
)

var (
	ErrBorradorInvalido     = errors.New("datos de borrador inválidos")
	ErrBorradorExistente    = errors.New("ya existe una inscripción activa para el estudiante, periodo y programa")
	ErrBorradorNoEncontrado = errors.New("borrador no encontrado para este estudiante")
	ErrBorradorCerrado      = errors.New("la versión no es un borrador editable")
	ErrRadicacionInvalida   = errors.New("la inscripción no cumple los requisitos de radicación")
)

func CrearBorradorInscripcionGrado(entrada models.CrearBorradorInscripcionGrado) (*models.BorradorInscripcionGrado, error) {
	if err := validarBorradorInscripcionGrado(&entrada); err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return nil, fmt.Errorf("iniciar transacción: %w", err)
	}
	rollback := func(cause error) (*models.BorradorInscripcionGrado, error) {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("%w; revertir transacción: %v", cause, rollbackErr)
		}
		return nil, cause
	}

	resultado := &models.BorradorInscripcionGrado{Soportes: []models.SoporteGrado{}}
	ahora := utils.HoraBogota()
	err = tx.Raw(`
		INSERT INTO paz_y_salvos.solicitud_grado
			(tercero_id, codigo_estudiante, periodo_id, programa_academico_id,
			 dependencia_oikos_id, calendario_evento_inscripcion_id,
			 calendario_evento_aprobacion_id, fecha_creacion, fecha_modificacion)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING *`,
		entrada.TerceroId, entrada.CodigoEstudiante, entrada.PeriodoId,
		entrada.ProgramaAcademicoId, entrada.DependenciaOikosId,
		entrada.CalendarioEventoInscripcionId, entrada.CalendarioEventoAprobacionId, ahora, ahora,
	).QueryRow(&resultado.Solicitud)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return rollback(ErrBorradorExistente)
		}
		return rollback(fmt.Errorf("crear solicitud de grado: %w", err))
	}

	err = tx.Raw(`
		INSERT INTO paz_y_salvos.formulario_solicitud_grado
			(solicitud_grado_id, version, contenido, fecha_creacion, fecha_modificacion)
		VALUES (?, 1, ?::jsonb, ?, ?)
		RETURNING *`, resultado.Solicitud.Id, string(entrada.Contenido), ahora, ahora,
	).QueryRow(&resultado.Formulario)
	if err != nil {
		return rollback(fmt.Errorf("crear formulario de grado: %w", err))
	}

	err = tx.Raw(`
		INSERT INTO paz_y_salvos.historial_solicitud_grado
			(solicitud_grado_id, formulario_solicitud_grado_id, tercero_id,
			 estado_solicitud_id, justificacion, fecha_creacion, fecha_modificacion)
		VALUES (?, ?, ?, ?, 'Borrador creado por el estudiante', ?, ?)
		RETURNING *`, resultado.Solicitud.Id, resultado.Formulario.Id,
		entrada.TerceroId, entrada.EstadoBorradorId, ahora, ahora,
	).QueryRow(&resultado.Historial)
	if err != nil {
		return rollback(fmt.Errorf("crear historial del borrador: %w", err))
	}
	if err := cargarContenidoFormulario(tx, &resultado.Formulario); err != nil {
		return rollback(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("confirmar creación del borrador: %w", err)
	}
	return resultado, nil
}

var (
	registroSNPValido = regexp.MustCompile(`^[[:alnum:]-]+$`)
	cedulaGradoValida = regexp.MustCompile(`^[1-9][0-9]{0,37}$`)
)

func contenidoRadicacionValido(contenido json.RawMessage) bool {
	if !contenidoBorradorValido(contenido) {
		return false
	}
	var formulario struct {
		TrabajoGrado               string `json:"trabajoGrado"`
		Director1                  string `json:"director1"`
		Director2                  string `json:"director2"`
		Modalidad                  string `json:"modalidad"`
		LugarExpedicionDocumentoId int64  `json:"lugarExpedicionDocumentoId"`
		LugarExpedicionDocumento   string `json:"lugarExpedicionDocumento"`
		NumeroActaSustentacion     string `json:"numeroActaSustentacion"`
		NumeroRegistroSNP          string `json:"numeroRegistroSnp"`
		TrabajaActualmente         *bool  `json:"trabajaActualmente"`
		Empresa                    string `json:"empresa"`
		DireccionEmpresa           string `json:"direccionEmpresa"`
		TelefonoEmpresa            string `json:"telefonoEmpresa"`
	}
	if json.Unmarshal(contenido, &formulario) != nil || strings.TrimSpace(formulario.TrabajoGrado) == "" ||
		!cedulaGradoValida.MatchString(strings.TrimSpace(formulario.Director1)) || strings.TrimSpace(formulario.Modalidad) == "" ||
		formulario.LugarExpedicionDocumentoId <= 0 || strings.TrimSpace(formulario.LugarExpedicionDocumento) == "" ||
		strings.TrimSpace(formulario.NumeroActaSustentacion) == "" ||
		!registroSNPValido.MatchString(strings.TrimSpace(formulario.NumeroRegistroSNP)) || formulario.TrabajaActualmente == nil {
		return false
	}
	director2 := strings.TrimSpace(formulario.Director2)
	if director2 != "" && (!cedulaGradoValida.MatchString(director2) || director2 == strings.TrimSpace(formulario.Director1)) {
		return false
	}
	return !*formulario.TrabajaActualmente || (strings.TrimSpace(formulario.Empresa) != "" &&
		strings.TrimSpace(formulario.DireccionEmpresa) != "" && strings.TrimSpace(formulario.TelefonoEmpresa) != "")
}

func validarBorradorInscripcionGrado(entrada *models.CrearBorradorInscripcionGrado) error {
	entrada.CodigoEstudiante = strings.TrimSpace(entrada.CodigoEstudiante)
	if entrada.TerceroId <= 0 || entrada.CodigoEstudiante == "" || entrada.PeriodoId <= 0 ||
		entrada.ProgramaAcademicoId <= 0 || entrada.DependenciaOikosId <= 0 ||
		entrada.CalendarioEventoInscripcionId <= 0 || entrada.CalendarioEventoAprobacionId <= 0 ||
		entrada.EstadoBorradorId <= 0 {
		return fmt.Errorf("%w: todos los identificadores y el código estudiantil son obligatorios", ErrBorradorInvalido)
	}
	if !contenidoBorradorValido(entrada.Contenido) {
		return fmt.Errorf("%w: Contenido debe ser un JSON válido distinto de null", ErrBorradorInvalido)
	}
	return nil
}

func contenidoBorradorValido(contenido json.RawMessage) bool {
	if !json.Valid(contenido) {
		return false
	}
	var objeto map[string]json.RawMessage
	return json.Unmarshal(contenido, &objeto) == nil && objeto != nil &&
		!bytes.Equal(bytes.TrimSpace(contenido), []byte("null"))
}

// ConsultarBorradorInscripcionGrado filtra por el propietario y los estados
// comprobados por el MID. Nunca devuelve expedientes de otro estudiante.
func ConsultarBorradorInscripcionGrado(terceroID, periodoID, programaID, estadoBorradorID, estadoRadicadaID int) (*models.BorradorInscripcionGrado, error) {
	if terceroID <= 0 || periodoID <= 0 || programaID <= 0 || estadoBorradorID <= 0 {
		return nil, ErrBorradorInvalido
	}
	o := orm.NewOrm()
	resultado := &models.BorradorInscripcionGrado{}
	if err := o.Raw(`SELECT * FROM paz_y_salvos.solicitud_grado
		WHERE tercero_id=? AND periodo_id=? AND programa_academico_id=? AND activo`,
		terceroID, periodoID, programaID).QueryRow(&resultado.Solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, ErrBorradorNoEncontrado
		}
		return nil, err
	}
	if err := cargarSolicitudConsultable(o, resultado, estadoBorradorID, estadoRadicadaID); err != nil {
		return nil, err
	}
	return resultado, nil
}

func ConsultarBorradorPorID(id, terceroID, estadoBorradorID, estadoRadicadaID int) (*models.BorradorInscripcionGrado, error) {
	if id <= 0 || terceroID <= 0 || estadoBorradorID <= 0 {
		return nil, ErrBorradorInvalido
	}
	o := orm.NewOrm()
	resultado := &models.BorradorInscripcionGrado{}
	if err := o.Raw(`SELECT * FROM paz_y_salvos.solicitud_grado
		WHERE id=? AND tercero_id=? AND activo`, id, terceroID).QueryRow(&resultado.Solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return nil, ErrBorradorNoEncontrado
		}
		return nil, err
	}
	if err := cargarSolicitudConsultable(o, resultado, estadoBorradorID, estadoRadicadaID); err != nil {
		return nil, err
	}
	return resultado, nil
}

func cargarBorrador(o orm.QueryExecutor, resultado *models.BorradorInscripcionGrado, estadoBorradorID int) error {
	if err := cargarSolicitud(o, resultado, estadoBorradorID); err != nil {
		return err
	}
	if resultado.Formulario.FechaRadicacion != nil {
		return ErrBorradorCerrado
	}
	return nil
}

func cargarSolicitudConsultable(o orm.QueryExecutor, resultado *models.BorradorInscripcionGrado, estadoBorradorID, estadoRadicadaID int) error {
	estados := []int{estadoBorradorID}
	if estadoRadicadaID > 0 && estadoRadicadaID != estadoBorradorID {
		estados = append(estados, estadoRadicadaID)
	}
	if err := cargarSolicitud(o, resultado, estados...); err != nil {
		return err
	}
	esRadicada := resultado.Historial.EstadoSolicitudId == estadoRadicadaID
	if (esRadicada && resultado.Formulario.FechaRadicacion == nil) || (!esRadicada && resultado.Formulario.FechaRadicacion != nil) {
		return ErrBorradorCerrado
	}
	return nil
}

func cargarSolicitud(o orm.QueryExecutor, resultado *models.BorradorInscripcionGrado, estadosPermitidos ...int) error {
	if err := o.Raw(`SELECT * FROM paz_y_salvos.formulario_solicitud_grado
		WHERE solicitud_grado_id=? AND activo
		ORDER BY version DESC LIMIT 1`, resultado.Solicitud.Id).QueryRow(&resultado.Formulario); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return ErrBorradorCerrado
		}
		return err
	}
	if err := o.Raw(`SELECT * FROM paz_y_salvos.historial_solicitud_grado
		WHERE solicitud_grado_id=? AND activo ORDER BY fecha_creacion DESC, id DESC LIMIT 1`,
		resultado.Solicitud.Id).QueryRow(&resultado.Historial); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return ErrBorradorCerrado
		}
		return err
	}
	estadoPermitido := false
	for _, estado := range estadosPermitidos {
		if estado > 0 && resultado.Historial.EstadoSolicitudId == estado {
			estadoPermitido = true
			break
		}
	}
	if !estadoPermitido || resultado.Historial.FormularioSolicitudGradoId == nil ||
		*resultado.Historial.FormularioSolicitudGradoId != resultado.Formulario.Id {
		return ErrBorradorCerrado
	}
	if err := cargarContenidoFormulario(o, &resultado.Formulario); err != nil {
		return err
	}
	resultado.Soportes = []models.SoporteGrado{}
	_, err := o.Raw(`SELECT * FROM paz_y_salvos.soporte_solicitud_grado
		WHERE solicitud_grado_id=? AND formulario_solicitud_grado_id=? AND activo ORDER BY tipo_documento_id`,
		resultado.Solicitud.Id, resultado.Formulario.Id).QueryRows(&resultado.Soportes)
	return err
}

func cargarContenidoFormulario(o orm.QueryExecutor, formulario *models.FormularioSolicitudGrado) error {
	var texto string
	if err := o.Raw(`SELECT contenido::text FROM paz_y_salvos.formulario_solicitud_grado WHERE id=?`,
		formulario.Id).QueryRow(&texto); err != nil {
		return err
	}
	formulario.Contenido = json.RawMessage(texto)
	return nil
}

// ActualizarBorradorInscripcionGrado serializa escrituras sobre la solicitud y
// exige propiedad, versión sin radicar y estado vigente borrador.
func ActualizarBorradorInscripcionGrado(id, terceroID, estadoBorradorID int, contenido json.RawMessage) (*models.BorradorInscripcionGrado, error) {
	if id <= 0 || terceroID <= 0 || estadoBorradorID <= 0 || !contenidoBorradorValido(contenido) {
		return nil, ErrBorradorInvalido
	}
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return nil, err
	}
	rollback := func(cause error) (*models.BorradorInscripcionGrado, error) {
		if e := tx.Rollback(); e != nil {
			return nil, fmt.Errorf("%w; revertir: %v", cause, e)
		}
		return nil, cause
	}
	resultado := &models.BorradorInscripcionGrado{}
	if err := tx.Raw(`SELECT * FROM paz_y_salvos.solicitud_grado WHERE id=? AND tercero_id=? AND activo FOR UPDATE`,
		id, terceroID).QueryRow(&resultado.Solicitud); err != nil {
		if errors.Is(err, orm.ErrNoRows) {
			return rollback(ErrBorradorNoEncontrado)
		}
		return rollback(err)
	}
	if err := cargarBorrador(tx, resultado, estadoBorradorID); err != nil {
		return rollback(err)
	}
	ahora := utils.HoraBogota()
	if err := tx.Raw(`UPDATE paz_y_salvos.formulario_solicitud_grado
		SET contenido=?::jsonb, fecha_modificacion=?
		WHERE id=? AND solicitud_grado_id=? AND fecha_radicacion IS NULL AND activo RETURNING *`,
		string(contenido), ahora, resultado.Formulario.Id, id).QueryRow(&resultado.Formulario); err != nil {
		return rollback(err)
	}
	if err := cargarContenidoFormulario(tx, &resultado.Formulario); err != nil {
		return rollback(err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return resultado, nil
}
