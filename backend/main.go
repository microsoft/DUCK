package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/labstack/echo/engine/standard"

	"github.com/Microsoft/DUCK/backend/ducklib"
	//Database lugin, change this if you have another Plugin/database
	_ "github.com/Microsoft/DUCK/backend/plugins/couchdb"
)

var (
	webDir      string
	jwtKey      string
	ruleBaseDir string
	goPath      = os.Getenv("GOPATH")
	startDir    = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/testdata.json")
)

func main() {

	flag.StringVar(&webDir, "webdir", "frontend/dist", "The root directory for serving web content")
	flag.StringVar(&jwtKey, "JWTSecret", "secret", "The secret used to sign the JWT")
	flag.StringVar(&ruleBaseDir, "rulebasedir", startDir, "The Directory to the Rulebases")

	flag.Parse()

	//create ComplianceCheckerPlugin

	//set routes
	e := ducklib.GetServer(webDir, []byte(jwtKey), ruleBaseDir)

	//start server
	e.Run(standard.New(":3000"))

}
