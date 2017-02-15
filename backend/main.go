package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/Microsoft/DUCK/backend/ducklib"
	"github.com/Microsoft/DUCK/backend/ducklib/config"

	//Database lugin, change this if you have another Plugin/database
	_ "github.com/Microsoft/DUCK/backend/plugins/couchdb"
)

//main function: loading config, starting server
func main() {

	//when panic occours print stack & then wait for user input to close
	//this is especially useful for windows users who start the compiled exe directly
	//and wont see the cause of an error when the program closes immediately afterwards
	defer func() {
		if r := recover(); r != nil {

			fmt.Printf("%s: %s", r, debug.Stack())

			fmt.Print("Press enter to exit ")
			var input string
			fmt.Scanln(&input)
			panic(r)

		}

	}()
	//when GOPATH is set, the person executing this program is probably also
	//writing code for it, it will thus also be in the goPath
	//when there is no GOPATH, the configuration.json should be right next to the executable
	confPath := "configuration.json"
	if goPath := os.Getenv("GOPATH"); goPath != "" {

		confPath = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/backend/configuration.json")
	}
	log.Printf("using configuration at %s", confPath)
	// create config
	conf := config.NewConfiguration(confPath)

	//set routes
	//	e := ducklib.GetServer(webDir, []byte(jwtKey), ruleBaseDir)
	e := ducklib.GetServer(conf)

	//start server
	e.Logger.Fatal(e.Start(":3000"))

}
