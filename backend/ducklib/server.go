package ducklib

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//structs

var datab *Database

//JWT contains the JWT secret
var JWT []byte

// Checker is a ComplianceCheckerPlugin
var checker *ComplianceCheckerPlugin

//GetServer returns Echo instance with predefined routes
func GetServer(webDir string, jwtKey []byte, ruleBaseDir string) *echo.Echo {

	datab = NewDatabase()
	err := datab.Init()

	JWT = jwtKey

	checker, err = MakeComplianceCheckerPlugin(ruleBaseDir)
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
	jwtMiddleware := middleware.JWT(jwtKey)
	users := api.Group("/users") //base URI

	//create a new user - JWT must not be required since during registration (when the user account is created) the user is not authenticated
	users.POST("", postUserHandler)
	users.DELETE("/:id", deleteUserHandler, jwtMiddleware) //delete a user
	users.PUT("/:id", putUserHandler, jwtMiddleware)       //update a user

	//data use statement document resources
	//documents := api.Group("/documents") //base URI
	documents := api.Group("/documents", jwtMiddleware) //base URI
	documents.POST("", postDocHandler)                  //create document
	documents.PUT("", putDocHandler)                    //update document
	documents.DELETE("/:docid", deleteDocHandler)       //delete document
	documents.GET("/:userid/summary", getDocSummaries)  //return document summaries for the author
	documents.GET("/:docid", getDocHandler)             //return document

	//ruleset resources
	rulebases := api.Group("/rulebases", jwtMiddleware) //base URI
	//rulesets.POST("/", postRsHandler)                                //create a ruleset
	//rulesets.DELETE("/:id", deleteRsHandler)                         //delete a ruleset
	//rulesets.PUT("/:setid", putRsHandler)                            //update a ruleset
	rulebases.PUT("/:baseid/documents", checkDocHandler)               //process provided document against ruleset
	rulebases.PUT("/:baseid/documents/:documentid", checkDocIDHandler) //process document against ruleset

	// serves the static files
	e.Static("/", webDir)

	return e

}
