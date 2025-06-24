package main

import (
	// <- nuevo

	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/lib/pq"
	"github.com/udistrital/auditoria"
	_ "github.com/udistrital/paz_y_salvos_crud/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	"github.com/udistrital/utils_oas/customerrorv2"
	"github.com/udistrital/utils_oas/ssm"
)

func init() {
	if beego.AppConfig.String("parameterStore") != "" {
		parameterStore := "/" + beego.AppConfig.String("parameterStore") +
			"/" + beego.AppConfig.String("appname") + "/db/"

		username, err := ssm.GetParameterFromParameterStore(parameterStore + "username")
		if err != nil {
			logs.Critical("Error retrieving username: %v", err)
		}

		err = beego.AppConfig.Set("PGuser", username)
		if err != nil {
			logs.Critical("Failed to set PGuser env var: %v", err)
		}

		password, err := ssm.GetParameterFromParameterStore(parameterStore + "password")
		if err != nil {
			logs.Critical("Error retrieving password: %v", err)
		}

		err = beego.AppConfig.Set("PGpass", password)
		if err != nil {
			logs.Critical("Failed to set PGpass: %v", err)
		}
	}

	orm.RegisterDataBase("default", "postgres", "postgres://"+beego.AppConfig.String("PGuser")+":"+url.QueryEscape(beego.AppConfig.String("PGpass"))+"@"+beego.AppConfig.String("PGhost")+":"+beego.AppConfig.String("PGport")+"/"+beego.AppConfig.String("PGdb")+"?sslmode=disable&search_path="+beego.AppConfig.String("PGschema"))
}

func main() {
	orm.RegisterDataBase("default", "postgres", "postgres://"+
		beego.AppConfig.String("PGuser")+":"+
		beego.AppConfig.String("PGpass")+"@"+
		beego.AppConfig.String("PGhost")+":"+
		beego.AppConfig.String("PGport")+"/"+
		beego.AppConfig.String("PGdb")+"?sslmode=disable&search_path="+
		beego.AppConfig.String("PGschema")+"")

	AllowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == "dev" {
		AllowedOrigins = []string{"*"}
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: AllowedOrigins,
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
}
