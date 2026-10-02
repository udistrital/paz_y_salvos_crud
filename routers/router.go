// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/udistrital/paz_y_salvos_crud/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/semaforo",
			beego.NSInclude(
				&controllers.SemaforoController{},
			),
		),
		beego.NSNamespace("/solicitud-grado",
			beego.NSRouter("/borrador/:id/soportes/:tipo", &controllers.InscripcionGradoController{}, "put:AsociarSoporte"),
			beego.NSRouter("/borrador/:id/radicar", &controllers.InscripcionGradoController{}, "post:Radicar"),
			beego.NSRouter("/borrador", &controllers.InscripcionGradoController{}, "post:CrearBorrador;get:ConsultarBorrador"),
			beego.NSRouter("/borrador/:id", &controllers.InscripcionGradoController{}, "get:ConsultarBorradorPorID;put:ActualizarBorrador"),
		),
	)
	beego.AddNamespace(ns)
}
