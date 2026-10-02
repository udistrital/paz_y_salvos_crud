package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/paz_y_salvos_crud/helpers"
	"github.com/udistrital/paz_y_salvos_crud/models"
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
	estadoID, e4 := c.GetInt("estado_borrador_id")
	estadoRadicadaID, e5 := parametroEnteroOpcional(c, "estado_radicada_id")
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	c.responderBorrador(services.ConsultarBorradorInscripcionGrado(terceroID, periodoID, programaID, estadoID, estadoRadicadaID))
}

func (c *InscripcionGradoController) ConsultarBorradorPorID() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	terceroID, e2 := c.GetInt("tercero_id")
	estadoID, e3 := c.GetInt("estado_borrador_id")
	estadoRadicadaID, e4 := parametroEnteroOpcional(c, "estado_radicada_id")
	if err != nil || e2 != nil || e3 != nil || e4 != nil {
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores inválidos")
		return
	}
	c.responderBorrador(services.ConsultarBorradorPorID(id, terceroID, estadoID, estadoRadicadaID))
}

func parametroEnteroOpcional(c *InscripcionGradoController, nombre string) (int, error) {
	valor := c.GetString(nombre)
	if valor == "" {
		return 0, nil
	}
	return strconv.Atoi(valor)
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
