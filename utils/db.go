package utils

import (
	"fmt"
	"net/url"
	"os"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/ssm"
)

func InitDB() (connStr string, username string, password string, err error) {
	username, password, err = getDBCredentials()
	if err != nil {
		return "", "", "", err
	}

	connStr = buildConnStr(username, password)
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
			return "", "", fmt.Errorf("failed to retrieve username from SSM: %v", err)
		}

		password, err = ssm.GetParameterFromParameterStore(basePath + "password")
		if err != nil {
			return "", "", fmt.Errorf("failed to retrieve password from SSM: %v", err)
		}
	} else {
		username = os.Getenv("PGUSER")
		password = os.Getenv("PGPASS")

		if username == "" || password == "" {
			return "", "", fmt.Errorf("missing PGUSER or PGPASS environment variables")
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
