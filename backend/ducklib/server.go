package ducklib

import (
	"log"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//structs

var datab *database

//JWT contains the JWT secret
var JWT []byte

// Checker is a ComplianceCheckerPlugin
var checker *ComplianceCheckerPlugin

//GetServer returns Echo instance with predefined routes
func GetServer(conf structs.Configuration) *echo.Echo {
	//webDir string, jwtKey []byte, ruleBaseDir string

	datab = NewDatabase(*conf.DBConfig)
	err := datab.Init()
	if err != nil {
		panic(err)
	}

	if conf.Loadtestdata {

		if err := FillTestdata(); err != nil {
			log.Printf("Error trying to load testdata: %s", err)
		}
	}

	JWT = []byte(conf.JwtKey)
	rbd := conf.RulebaseDir

	log.Printf("Rulebase directory: " + rbd)

	checker, err = MakeComplianceCheckerPlugin(rbd)
	if err != nil {
		panic(err)
	}
	err = checker.Intialize()
	if err != nil {
		panic(err)
	}

	//New echo instance
	e := echo.New()

	//set used Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	//Logger Config
	LoggerConfig := middleware.LoggerConfig{Format: `{"time":"${time_rfc3339}",` +
		`"method":"${method}","uri":"${uri}","status":${status}, ` +
		`"Bytes received":${rx_bytes},"Bytes sent":${tx_bytes}}` + "\n",
	}
	e.Use(middleware.LoggerWithConfig(LoggerConfig))
	e.Use(middleware.Recover())

	e.POST("/login", loginHandler)
	e.GET("/loadtestdata", testdataHandler)
	//create sub-router for api functions
	api := e.Group("/v1")

	////User resources
	jwtMiddleware := middleware.JWT(JWT)
	users := api.Group("/users") //base URI

	//create a new user - JWT must not be required since during registration (when the user account is created) the user is not authenticated
	users.POST("", postUserHandler)
	users.DELETE("/:id", deleteUserHandler, jwtMiddleware)                      //delete a user
	users.PUT("/", putUserHandler, jwtMiddleware)                               //update a user
	users.GET("/:id/dictionary", getUserDictHandler, jwtMiddleware)             //get a users dictonary
	users.PUT("/:id/dictionary", putUserDictHandler, jwtMiddleware)             //update a users dictonary
	users.GET("/:id/dictionary/:code", getDictItemHandler, jwtMiddleware)       //get a dictonary entry
	users.PUT("/:id/dictionary/:code", putDictItemHandler, jwtMiddleware)       //update a dictonary entry
	users.DELETE("/:id/dictionary/:code", deleteDictItemHandler, jwtMiddleware) //delete a dictonary entry

	//data use statement document resources
	documents := api.Group("/documents", jwtMiddleware)   //base URI
	documents.POST("", postDocHandler)                    //create document
	documents.PUT("", putDocHandler)                      //update document
	documents.DELETE("/:docid", deleteDocHandler)         //delete document
	documents.GET("/:userid/summary", getDocSummaries)    //return document summaries for the author
	documents.GET("/:docid", getDocHandler)               //return document
	documents.POST("/copy/:docid", copyStatementsHandler) //copies the statements from an existing Document to a new one

	//rulebase resources
	rulebases := api.Group("/rulebases", jwtMiddleware) //base URI
	rulebases.GET("", getRulebasesHandler)              //Returns a dictionary with all available Rulebases
	//rulebases.POST("/", postRsHandler)                                //create a rulebase
	//rulebases.DELETE("/:id", deleteRsHandler)                         //delete a rulebase
	//rulebases.PUT("/:setid", putRsHandler)                            //update a rulebase
	rulebases.PUT("/:baseid/documents", checkDocHandler)               //process provided document against rulebase
	rulebases.PUT("/:baseid/documents/:documentid", checkDocIDHandler) //process document against rulebase

	// serves the static files
	wbd := conf.WebDir

	log.Printf("Web directory: " + wbd)
	e.Static("/", wbd)

	return e

}
