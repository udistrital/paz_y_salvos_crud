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

// SemaforoController operations for Semaforo
type SemaforoController struct {
	beego.Controller
}

// URLMapping ...
func (c *SemaforoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Patch", c.Patch)
}

// Post ...
// @Title Post
// @Description create Semaforo
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

// GetOne ...
// @Title Get One
// @Description get Semaforo by id
// @Param	id		path 	string	true		"The key for staticblock"
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

// GetAll ...
// @Title Get All
// @Description get Semaforo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
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

// Put ...
// @Title Put
// @Description update the Semaforo
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
	v.Id = id // Asegura que el id es el correcto

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

// Delete ...
// @Title Delete
// @Description delete the Semaforo
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

// Patch ...
// @Title Patch
// @Description update partial fields of Semaforo
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

		// Si es un error de validación de reglas de negocio, retornar 409 Conflict
		if strings.Contains(errorMsg, "ORC") ||
			strings.Contains(errorMsg, "dependencias") ||
			strings.Contains(errorMsg, "activo") {
			helpers.RenderError(&c.Controller, http.StatusConflict, errorMsg)
			return
		}

		// Para otros errores, retornar 400
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
