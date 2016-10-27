package main

import (
	"os"
	"path/filepath"

	"github.com/Microsoft/DUCK/backend/ducklib"
	"github.com/Microsoft/DUCK/backend/ducklib/config"
	"github.com/labstack/echo/engine/standard"

	//Database lugin, change this if you have another Plugin/database
	_ "github.com/Microsoft/DUCK/backend/plugins/couchdb"
)

func main() {

	goPath := os.Getenv("GOPATH")
	confPath := "configuration.json"
	if goPath != "" {
		confPath = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/configuration.json")
	}
	// create config
	conf := config.NewConfiguration(confPath)

	//set routes
	//	e := ducklib.GetServer(webDir, []byte(jwtKey), ruleBaseDir)
	e := ducklib.GetServer(conf)

	//start server
	e.Run(standard.New(":3000"))

}
