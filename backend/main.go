package main

import (
	"flag"

	"github.com/labstack/echo/engine/standard"

	"github.com/Metaform/duck/backend/ducklib"
	//Database lugin, change this if you have another Plugin/database
	_ "github.com/Metaform/duck/backend/plugins/couchbase"
)

var (
	webDir string
	jwtKey = []byte("secret")
)

//loads config & checks if db has to be setup
func init() {
	flag.StringVar(&webDir, "webdir", "frontend/dist", "The root directory for serving web content")

}

func main() {
	flag.Parse()
	//fmt.Println("Web root: " + webDir)

	//set routes
	e := ducklib.GetServer(webDir, jwtKey)

	//start server
	e.Run(standard.New(":3000"))

}
