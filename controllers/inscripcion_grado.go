package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/paz_y_salvos_crud/helpers"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/repositories"
	"github.com/udistrital/paz_y_salvos_crud/services"
)

type InscripcionGradoController struct {
	beego.Controller
}

func (c *InscripcionGradoController) leerJSON(destino interface{}) bool {
	decoder := json.NewDecoder(bytes.NewReader(c.Ctx.Input.RequestBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destino); err != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "JSON inválido o campos no reconocidos")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Se requiere un único objeto JSON")
		return false
	}
	return true
}

func (c *InscripcionGradoController) CrearBorrador() {
	var entrada models.CrearBorradorInscripcionGrado
	if !c.leerJSON(&entrada) {
		return
	}

	borrador, err := services.CrearBorradorInscripcionGrado(entrada)
	if err != nil {
		logs.Error(err)
		switch {
		case errors.Is(err, services.ErrBorradorInvalido):
			helpers.RenderError(&c.Controller, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrBorradorExistente):
			helpers.RenderError(&c.Controller, http.StatusConflict, err.Error())
		default:
			helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo crear el borrador de inscripción a grado")
		}
		return
	}

	helpers.RenderResponse(&c.Controller, models.APIResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Borrador de inscripción a grado creado correctamente",
		Data:    borrador,
	})
}

func (c *InscripcionGradoController) ConsultarBorrador() {
	terceroID, e1 := c.GetInt("tercero_id")
	periodoID, e2 := c.GetInt("periodo_id")
	programaID, e3 := c.GetInt("programa_id")
	estados, e4 := idsConsultaRevision(c.GetString("estados"))
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	c.responderBorrador(services.ConsultarBorradorInscripcionGrado(terceroID, periodoID, programaID, estados...))
}

func (c *InscripcionGradoController) ConsultarBorradorPorID() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	terceroID, e2 := c.GetInt("tercero_id")
	estados, e3 := idsConsultaRevision(c.GetString("estados"))
	if err != nil || e2 != nil || e3 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	c.responderBorrador(services.ConsultarBorradorPorID(id, terceroID, estados...))
}

func (c *InscripcionGradoController) ActualizarBorrador() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	terceroID, e2 := c.GetInt("tercero_id")
	estadoID, e3 := c.GetInt("estado_borrador_id")
	if err != nil || e2 != nil || e3 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	var entrada models.ActualizarBorradorInscripcionGrado
	if !c.leerJSON(&entrada) {
		return
	}
	c.responderBorrador(services.ActualizarBorradorInscripcionGrado(id, terceroID, estadoID, entrada.Contenido))
}

func (c *InscripcionGradoController) AsociarSoporte() {
	id, e1 := strconv.Atoi(c.Ctx.Input.Param(":id"))
	tipo, e2 := strconv.Atoi(c.Ctx.Input.Param(":tipo"))
	tercero, e3 := c.GetInt("tercero_id")
	estado, e4 := c.GetInt("estado_borrador_id")
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		helpers.RenderError(&c.Controller, 400, "Identificadores inválidos")
		return
	}
	var entrada models.AsociarSoporteGrado
	if !c.leerJSON(&entrada) {
		return
	}
	c.responderBorrador(services.AsociarSoporteBorrador(id, tercero, estado, tipo, entrada))
}

func (c *InscripcionGradoController) EliminarSoporte() {
	id, e1 := strconv.Atoi(c.Ctx.Input.Param(":id"))
	tipo, e2 := strconv.Atoi(c.Ctx.Input.Param(":tipo"))
	tercero, e3 := c.GetInt("tercero_id")
	estado, e4 := c.GetInt("estado_borrador_id")
	formulario, e5 := c.GetInt("formulario_id")
	soporte, e6 := c.GetInt("soporte_actual_id")
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil || e6 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	c.responderBorrador(services.EliminarSoporteBorrador(id, tercero, estado, tipo, formulario, soporte))
}

func (c *InscripcionGradoController) Radicar() {
	id, e1 := strconv.Atoi(c.Ctx.Input.Param(":id"))
	tercero, e2 := c.GetInt("tercero_id")
	if e1 != nil || e2 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	var entrada models.RadicarInscripcionGrado
	if !c.leerJSON(&entrada) {
		return
	}
	resultado, err := services.RadicarInscripcionGrado(id, tercero, entrada)
	if err != nil {
		logs.Error(err)
		switch {
		case errors.Is(err, services.ErrRadicacionInvalida):
			helpers.RenderError(&c.Controller, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrBorradorNoEncontrado):
			helpers.RenderError(&c.Controller, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrBorradorCerrado):
			helpers.RenderError(&c.Controller, http.StatusConflict, err.Error())
		default:
			helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo radicar la inscripción")
		}
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Inscripción radicada", Data: resultado})
}

func (c *InscripcionGradoController) Subsanar() {
	id, e1 := strconv.Atoi(c.Ctx.Input.Param(":id"))
	tercero, e2 := c.GetInt("tercero_id")
	if e1 != nil || e2 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	var entrada models.SubsanarInscripcionGrado
	if !c.leerJSON(&entrada) {
		return
	}
	c.responderBorrador(services.SubsanarInscripcionGrado(id, tercero, entrada))
}

func idsConsultaRevision(valor string) ([]int, error) {
	partes := strings.Split(valor, ",")
	ids := make([]int, 0, len(partes))
	vistos := make(map[int]bool, len(partes))
	for _, parte := range partes {
		id, err := strconv.Atoi(strings.TrimSpace(parte))
		if err != nil || id <= 0 || vistos[id] {
			return nil, services.ErrRevisionDocumentalInvalida
		}
		vistos[id] = true
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, services.ErrRevisionDocumentalInvalida
	}
	return ids, nil
}

func (c *InscripcionGradoController) ListarRevisionDocumental() {
	dependencia, errDependencia := c.GetInt("dependencia_id")
	estados, errEstados := idsConsultaRevision(c.GetString("estados"))
	limit, errLimit := c.GetInt("limit", 20)
	offset, errOffset := c.GetInt("offset", 0)
	periodo, errPeriodo := c.GetInt("periodo_id", 0)
	programa, errPrograma := c.GetInt("programa_id", 0)
	if errDependencia != nil || errEstados != nil || errLimit != nil || errOffset != nil || errPeriodo != nil || errPrograma != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Filtros de revisión inválidos")
		return
	}
	resultado, err := services.ListarSolicitudesRevisionGrado(repositories.NuevaRevisionDocumentalORM(), dependencia, estados, limit, offset, periodo, programa)
	if err != nil {
		c.responderRevision(nil, err)
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Solicitudes para revisión", Data: resultado})
}

func (c *InscripcionGradoController) ConsultarRevisionDocumental() {
	id, errID := strconv.Atoi(c.Ctx.Input.Param(":id"))
	dependencia, errDependencia := c.GetInt("dependencia_id")
	estados, errEstados := idsConsultaRevision(c.GetString("estados"))
	if errID != nil || errDependencia != nil || errEstados != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores de revisión inválidos")
		return
	}
	c.responderRevision(services.ConsultarSolicitudRevisionGrado(repositories.NuevaRevisionDocumentalORM(), id, dependencia, estados))
}

func (c *InscripcionGradoController) RevisarDocumentacion() {
	id, errID := strconv.Atoi(c.Ctx.Input.Param(":id"))
	dependencia, errDependencia := c.GetInt("dependencia_id")
	if errID != nil || errDependencia != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores de revisión inválidos")
		return
	}
	var entrada models.RevisarDocumentacionGrado
	if !c.leerJSON(&entrada) {
		return
	}
	c.responderRevision(services.RevisarDocumentacionGrado(repositories.NuevaRevisionDocumentalORM(), id, dependencia, entrada))
}

func (c *InscripcionGradoController) ConsultarPazSalvos() {
	id, errID := strconv.Atoi(c.Ctx.Input.Param(":id"))
	estado, errEstado := c.GetInt("estado_documentacion_aprobada_id")
	tipos, errTipos := idsConsultaRevision(c.GetString("tipos"))
	if errID != nil || errEstado != nil || errTipos != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores de Paz y Salvo inválidos")
		return
	}
	c.responderPazSalvo(services.ConsultarPazSalvosGrado(repositories.NuevoPazSalvoGradoORM(), id, estado, tipos))
}

func (c *InscripcionGradoController) ListarPazSalvos() {
	estado, errEstado := c.GetInt("estado_documentacion_aprobada_id")
	tipos, errTipos := idsConsultaRevision(c.GetString("tipos"))
	limit, errLimit := c.GetInt("limit", 20)
	offset, errOffset := c.GetInt("offset", 0)
	periodo, errPeriodo := c.GetInt("periodo_id", 0)
	programa, errPrograma := c.GetInt("programa_id", 0)
	tercero, errTercero := c.GetInt("tercero_id", 0)
	var dependencias []int
	var errDependencias error
	if strings.TrimSpace(c.GetString("dependencias")) != "" {
		dependencias, errDependencias = idsConsultaRevision(c.GetString("dependencias"))
	}
	if errEstado != nil || errTipos != nil || errLimit != nil || errOffset != nil || errPeriodo != nil || errPrograma != nil || errTercero != nil || errDependencias != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Filtros de Paz y Salvos inválidos")
		return
	}
	resultado, err := services.ListarPazSalvosGrado(repositories.NuevoPazSalvoGradoORM(), estado, tipos, dependencias, limit, offset, periodo, programa, tercero, c.GetString("codigo"))
	if err != nil {
		if errors.Is(err, services.ErrPazSalvoGradoInvalido) {
			helpers.RenderError(&c.Controller, http.StatusBadRequest, err.Error())
			return
		}
		helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo consultar la bandeja de Paz y Salvos")
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Bandeja de Paz y Salvos", Data: resultado})
}

func (c *InscripcionGradoController) DecidirPazSalvo() {
	id, errID := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if errID != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Solicitud de Paz y Salvo inválida")
		return
	}
	var entrada models.DecidirPazSalvoGrado
	if !c.leerJSON(&entrada) {
		return
	}
	c.responderPazSalvo(services.DecidirPazSalvoGrado(repositories.NuevoPazSalvoGradoORM(), id, entrada))
}

func (c *InscripcionGradoController) responderPazSalvo(resultado *models.PazSalvosSolicitudGrado, err error) {
	if err != nil {
		logs.Error(err)
		switch {
		case errors.Is(err, services.ErrPazSalvoGradoInvalido):
			helpers.RenderError(&c.Controller, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrPazSalvoGradoNoEncontrado):
			helpers.RenderError(&c.Controller, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrPazSalvoGradoConflicto):
			helpers.RenderError(&c.Controller, http.StatusConflict, err.Error())
		default:
			helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo procesar el Paz y Salvo")
		}
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{
		Success: true, Status: http.StatusOK, Message: "Paz y Salvo procesado", Data: resultado,
	})
}

func (c *InscripcionGradoController) responderRevision(resultado *models.BorradorInscripcionGrado, err error) {
	if err != nil {
		logs.Error(err)
		switch {
		case errors.Is(err, services.ErrRevisionDocumentalInvalida):
			helpers.RenderError(&c.Controller, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrSolicitudFueraAlcance), errors.Is(err, services.ErrBorradorNoEncontrado):
			helpers.RenderError(&c.Controller, http.StatusNotFound, "Solicitud no encontrada en el alcance autorizado")
		case errors.Is(err, services.ErrBorradorCerrado):
			helpers.RenderError(&c.Controller, http.StatusConflict, "La solicitud cambió de estado o versión")
		default:
			helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo procesar la revisión documental")
		}
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Revisión documental registrada", Data: resultado})
}

func (c *InscripcionGradoController) responderBorrador(resultado *models.BorradorInscripcionGrado, err error) {
	if err != nil {
		logs.Error(err)
		switch {
		case errors.Is(err, services.ErrBorradorInvalido):
			helpers.RenderError(&c.Controller, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrBorradorNoEncontrado):
			helpers.RenderError(&c.Controller, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrBorradorCerrado):
			helpers.RenderError(&c.Controller, http.StatusConflict, err.Error())
		default:
			helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo consultar o guardar el borrador")
		}
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Consulta exitosa", Data: resultado})
}
