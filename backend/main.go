package main

import (
	"os"
	"path/filepath"

	"github.com/Microsoft/DUCK/backend/ducklib"
	"github.com/labstack/echo/engine/standard"
)

//Database lugin, change this if you have another Plugin/database
import _ "github.com/Microsoft/DUCK/backend/plugins/couchdb"

var (
	goPath = os.Getenv("GOPATH")
)

func main() {

	// create config
	conf := ducklib.NewConfiguration(filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/configuration.json"))

	//set routes
	//	e := ducklib.GetServer(webDir, []byte(jwtKey), ruleBaseDir)
	e := ducklib.GetServer(conf, goPath)

	//start server
	e.Run(standard.New(":3000"))

}
