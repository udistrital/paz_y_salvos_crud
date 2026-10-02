package helpers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/paz_y_salvos_crud/models"
)

func RenderResponse(c *beego.Controller, response models.APIResponse) {
	c.Ctx.Output.SetStatus(response.Status)
	c.Data["json"] = response
	if err := c.ServeJSON(); err != nil {
		logs.Error("error al serializar la respuesta: %v", err)
	}
}

func RenderError(c *beego.Controller, status int, message string) {
	RenderResponse(c, models.APIResponse{
		Success: false,
		Status:  status,
		Message: message,
		Data:    nil,
	})
}
