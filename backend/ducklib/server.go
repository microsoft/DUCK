package ducklib

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//structs

var datab *Database

//JWT contains the JWT secret
var JWT []byte

//GetServer returns Echo instance with predefined routes
func GetServer(webDir string, jwtKey []byte) *echo.Echo {

	datab = NewDatabase()
	datab.Init()

	JWT = jwtKey
	//New echo instance
	e := echo.New()

	//set used Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	//Logger Config
	LoggerConfig := middleware.LoggerConfig{Format: `{"time":"${time_rfc3339}",` +
		`"method":"${method}","uri":"${uri}","status":${status}, ` +
		`"latency":"${latency_human}","Bytes received":${rx_bytes},` +
		`"Bytes sent":${tx_bytes}}` + "\n",
	}
	e.Use(middleware.LoggerWithConfig(LoggerConfig))
	e.Use(middleware.Recover())

	
	e.POST("/login", loginHandler)
	e.GET("/loadtestdata", testdataHandler)
	//create sub-router for api functions
	api := e.Group("/v1")

	////User resources
	users := api.Group("/users", middleware.JWT(jwtKey)) //base URI

	users.POST("/", postUserHandler)        //create a new user
	users.DELETE("/:id", deleteUserHandler) //delete a user
	users.PUT("/:id", helloHandler)         //update a user

	//data use statement document resources
	//documents := api.Group("/documents") //base URI
	documents := api.Group("/documents", middleware.JWT(jwtKey)) //base URI
	documents.POST("", postDocHandler)                           //create document
	documents.PUT("", putDocHandler)                             //update document
	documents.DELETE("/:docid", deleteDocHandler)                //delete document
	documents.GET("/:userid/summary", getDocSummaries)           //return document summaries for the author
	documents.GET("/:docid", getDocHandler)                      //return document

	//ruleset resources
	rulesets := api.Group("/rulesets", middleware.JWT(jwtKey))       //base URI
	rulesets.POST("/", postRsHandler)                                //create a ruleset
	rulesets.DELETE("/:id", deleteRsHandler)                         //delete a ruleset
	rulesets.PUT("/:setid", putRsHandler)                            //update a ruleset
	rulesets.PUT("/:setid/documents", checkDocHandler)               //process provided document against ruleset
	rulesets.PUT("/:setid/documents/:documentid", checkDocIDHandler) //process document against ruleset

	// serves the static files
	e.Static("/", webDir)

	return e

}
