package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/udistrital/paz_y_salvos_crud/helpers"
	"github.com/udistrital/paz_y_salvos_crud/models"
	"github.com/udistrital/paz_y_salvos_crud/services"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// SemaforoController administra los registros del semáforo académico.
type SemaforoController struct {
	beego.Controller
}

// URLMapping vincula las operaciones del semáforo con los métodos HTTP.
func (c *SemaforoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Patch", c.Patch)
}

// Post crea un registro de semáforo a partir del cuerpo JSON.
// @Title Post
// @Description Crea un registro del semáforo
// @Param	body		body 	models.Semaforo	true		"body for Semaforo content"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @router / [post]
func (c *SemaforoController) Post() {
	var v models.Semaforo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		return
	}

	if _, err := models.AddSemaforo(&v); err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "No se pudo crear el semáforo con los datos enviados")
		return
	}

	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusCreated, Message: "Registro creado correctamente", Data: v})
}

// GetOne consulta un registro de semáforo por ID.
// @Title Get One
// @Description Consulta un registro del semáforo por ID
// @Param	id		path 	string	true		"ID del semáforo"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @router /:id [get]
func (c *SemaforoController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "El identificador debe ser un número entero")
		return
	}

	v, err := models.GetSemaforoById(id)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusNotFound, "No se encontró el semáforo solicitado")
		return
	}

	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Consulta exitosa", Data: v})
}

// GetAll lista registros de semáforo con filtros, orden y paginación.
// @Title Get All
// @Description Lista registros del semáforo con filtros opcionales
// @Param	query	query	string	false	"Filtros por campo en formato campo:valor"
// @Param	fields	query	string	false	"Campos incluidos en la respuesta, separados por coma"
// @Param	sortby	query	string	false	"Campos por los que se ordenan los resultados"
// @Param	order	query	string	false	"Dirección asc o desc por campo, o una sola para todos"
// @Param	limit	query	string	false	"Cantidad máxima de registros por página"
// @Param	offset	query	string	false	"Registros omitidos al inicio de la página"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @router / [get]
func (c *SemaforoController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				helpers.RenderError(&c.Controller, http.StatusBadRequest, "El filtro de consulta debe tener el formato clave:valor")
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllSemaforo(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusNotFound, "No se encontraron semáforos para los parámetros enviados")
		return
	}

	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Consulta exitosa", Data: l})
}

// Put reemplaza un registro del semáforo por ID y devuelve los datos persistidos.
// @Title Put
// @Description Actualiza un registro del semáforo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Semaforo	true		"body for Semaforo content"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @router /:id [put]
func (c *SemaforoController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "El identificador debe ser un número entero")
		return
	}

	var v models.Semaforo
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		return
	}
	v.Id = id // El ID de la ruta prevalece sobre el ID recibido en el cuerpo.

	if err := models.UpdateSemaforoById(&v); err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "No se pudo actualizar el semáforo con los datos enviados")
		return
	}

	// Consulta el registro actualizado para obtener los valores reales de la BD
	updated, err := models.GetSemaforoById(id)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusNotFound, "No se encontró el semáforo actualizado")
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Registro actualizado correctamente", Data: updated})
}

// Delete elimina el registro del semáforo indicado por ID.
// @Title Delete
// @Description Elimina un registro del semáforo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @router /:id [delete]
func (c *SemaforoController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "El identificador debe ser un número entero")
		return
	}

	if err := models.DeleteSemaforo(id); err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusNotFound, "No se encontró el semáforo que se desea eliminar")
		return
	}

	data := models.DeletedSemaforoResponse{Id: id}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Registro eliminado correctamente", Data: data})
}

// Patch aplica una actualización parcial y valida las reglas del semáforo.
// @Title Patch
// @Description Actualiza campos del semáforo y valida las reglas del negocio
// @Param   id      path    string                 true        "The id you want to patch"
// @Param   body    body    models.SemaforoPatch true        "Fields to update"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @router /:id [patch]
func (c *SemaforoController) Patch() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "El identificador debe ser un número entero")
		return
	}

	var patch models.SemaforoPatch
	decoder := json.NewDecoder(bytes.NewReader(c.Ctx.Input.RequestBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&patch); err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		return
	}

	if err := services.PatchSemaforoWithValidation(id, patch); err != nil {
		logs.Error(err)
		errorMsg := err.Error()

		// Los conflictos de estado, dependencia y actividad reciben HTTP 409.
		if strings.Contains(errorMsg, "ORC") ||
			strings.Contains(errorMsg, "dependencias") ||
			strings.Contains(errorMsg, "activo") {
			helpers.RenderError(&c.Controller, http.StatusConflict, errorMsg)
			return
		}

		// Las demás entradas inválidas reciben HTTP 400.
		helpers.RenderError(&c.Controller, http.StatusBadRequest, "No se pudo actualizar parcialmente el semáforo con los datos enviados")
		return
	}

	updated, err := models.GetSemaforoById(id)
	if err != nil {
		logs.Error(err)
		helpers.RenderError(&c.Controller, http.StatusNotFound, "No se encontró el semáforo actualizado")
		return
	}
	helpers.RenderResponse(&c.Controller, models.APIResponse{Success: true, Status: http.StatusOK, Message: "Registro actualizado parcialmente", Data: updated})
}
