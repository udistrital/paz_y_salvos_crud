package main

import (
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/lib/pq"
	_ "github.com/udistrital/paz_y_salvos_crud/routers"
	"github.com/udistrital/paz_y_salvos_crud/utils"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	"github.com/udistrital/utils_oas/auditoria"
	"github.com/udistrital/utils_oas/customerrorv2"
	"github.com/udistrital/utils_oas/xray"
)

func main() {
	connStr, username, password, err := utils.InitDB()
	if err != nil {
		logs.Critical("Error: %v", err)
		os.Exit(1)
	}

	_ = beego.AppConfig.Set("PGuser", username)
	_ = beego.AppConfig.Set("PGpass", password)

	if err := orm.RegisterDataBase("default", "postgres", connStr); err != nil {
		logs.Critical("Failed to register the database with ORM: %v", err)
		os.Exit(1)
	}
	logs.Info("Database registered successfully with ORM")

	// CORS y configuraciones de entorno
	AllowedOrigins := []string{"*.udistrital.edu.co"}
	if beego.BConfig.RunMode == "dev" {
		AllowedOrigins = []string{"*"}
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     AllowedOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "x-requested-with", "content-type", "accept", "authorization", "x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	xray.InitXRay()
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	apistatus.Init()
	auditoria.InitMiddleware()
	beego.Run()
}
