package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Microsoft/DUCK/backend/ducklib"
	"github.com/labstack/echo/engine/standard"
)

//Database lugin, change this if you have another Plugin/database
import _ "github.com/Microsoft/DUCK/backend/plugins/couchdb"

var (
	webDir      string
	jwtKey      string
	ruleBaseDir string
	goPath      = os.Getenv("GOPATH")
	startDir    = "/src/github.com/Microsoft/DUCK/RuleBases"
)

func main() {

	flag.StringVar(&webDir, "webdir", "frontend/dist", "The root directory for serving web content")
	flag.StringVar(&jwtKey, "JWTSecret", "secret", "The secret used to sign the JWT")
	flag.StringVar(&ruleBaseDir, "rulebasedir", startDir, "The Directory to the Rulebases")

	conf := ducklib.NewConfiguration(filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/configuration.json"))

	//set routes
	//	e := ducklib.GetServer(webDir, []byte(jwtKey), ruleBaseDir)
	e := ducklib.GetServer(conf, goPath)

	//start server
	e.Run(standard.New(":3000"))

}
