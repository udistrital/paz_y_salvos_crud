package utils

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/lib/pq"
	"github.com/udistrital/utils_oas/ssm"
)

func InitDB() (connStr string, username string, password string, err error) {
	username, password, err = getDBCredentials()
	if err != nil {
		return "", "", "", err
	}

	connStr = buildConnStr(username, password)

	err = validateDBConnection(connStr)
	if err != nil {
		return "", "", "", err
	}
	logs.Info("Database connection established successfully")
	return connStr, username, password, nil
}

func getDBCredentials() (string, string, error) {
	var username, password string
	var err error

	parameterStore := beego.AppConfig.String("parameterStore")
	appName := beego.AppConfig.String("appname")

	if parameterStore != "" {
		basePath := fmt.Sprintf("/%s/%s/db/", parameterStore, appName)

		username, err = ssm.GetParameterFromParameterStore(basePath + "username")
		if err != nil {
			return "", "", fmt.Errorf("error retrieving username: %v", err)
		}

		password, err = ssm.GetParameterFromParameterStore(basePath + "password")
		if err != nil {
			return "", "", fmt.Errorf("error retrieving password: %v", err)
		}
	} else {
		username = os.Getenv("PGUSER")
		password = os.Getenv("PGPASS")

		if username == "" || password == "" {
			return "", "", fmt.Errorf("error retrieving PGUSER and PGPASS from env variables")
		}
	}
	return username, password, nil
}

func buildConnStr(user, pass string) string {
	return "postgres://" + user + ":" + url.QueryEscape(pass) + "@" +
		beego.AppConfig.String("PGhost") + ":" +
		beego.AppConfig.String("PGport") + "/" +
		beego.AppConfig.String("PGdb") + "?sslmode=disable&search_path=" +
		beego.AppConfig.String("PGschema")
}

func validateDBConnection(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open connection to the database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("unable to establish connection to the database (ping failed): %v", err)
	}
	return nil
}
