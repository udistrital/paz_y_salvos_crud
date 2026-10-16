package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/udistrital/paz_y_salvos_crud/helpers"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/services"
)

func (c *InscripcionGradoController) ambitoVersion() (int, int, int, error) {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		return 0, 0, 0, services.ErrBorradorInvalido
	}
	tercero, e1 := c.GetInt("tercero_id", 0)
	dependencia, e2 := c.GetInt("dependencia_id", 0)
	if e1 != nil || e2 != nil || id <= 0 || (tercero <= 0) == (dependencia <= 0) {
		return 0, 0, 0, services.ErrBorradorInvalido
	}
	return id, tercero, dependencia, nil
}

func (c *InscripcionGradoController) responderVersion(data interface{}, err error) {
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBorradorInvalido):
			helpers.RenderError(&c.Controller, http.StatusBadRequest, "Identificadores de versión inválidos")
		case errors.Is(err, services.ErrBorradorNoEncontrado):
			helpers.RenderError(&c.Controller, http.StatusNotFound, "Entrega no encontrada")
		case errors.Is(err, services.ErrBorradorCerrado):
			helpers.RenderError(&c.Controller, http.StatusConflict, "Entrega no verificable")
		default:
			helpers.RenderError(&c.Controller, http.StatusInternalServerError, "No se pudo consultar la entrega")
		}
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Consulta exitosa", Data: data})
}

func (c *InscripcionGradoController) ListarVersiones() {
	id, tercero, dependencia, err := c.ambitoVersion()
	if err != nil {
		c.responderVersion(nil, err)
		return
	}
	versiones, err := services.ListarVersionesGrado(id, tercero, dependencia)
	c.responderVersion(versiones, err)
}

func (c *InscripcionGradoController) ConsultarVersion() {
	id, tercero, dependencia, err := c.ambitoVersion()
	if err != nil {
		c.responderVersion(nil, err)
		return
	}
	version, err := strconv.Atoi(c.Ctx.Input.Param(":numero"))
	if err != nil || version <= 0 {
		c.responderVersion(nil, services.ErrBorradorInvalido)
		return
	}
	detalle, err := services.ConsultarVersionGrado(id, version, tercero, dependencia)
	c.responderVersion(detalle, err)
}
