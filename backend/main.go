package main

import (
	"flag"

	"github.com/labstack/echo/engine/standard"

	"github.com/Microsoft/DUCK/backend/ducklib"
	//Database lugin, change this if you have another Plugin/database
	_ "github.com/Microsoft/DUCK/backend/plugins/couchbase"
)

var (
	webDir string
	jwtKey string
)

func main() {

	flag.StringVar(&webDir, "webdir", "frontend/dist", "The root directory for serving web content")
	flag.StringVar(&jwtKey, "JWTSecret", "secret", "The secret used to sign the JWT")

	flag.Parse()
	//fmt.Println("Web root: " + webDir)

	//set routes
	e := ducklib.GetServer(webDir, []byte(jwtKey))

	//start server
	e.Run(standard.New(":3000"))

}
