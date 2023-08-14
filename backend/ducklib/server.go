// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package ducklib

import (
	"log"

	"github.com/Microsoft/DUCK/backend/ducklib/carneades"
	"github.com/Microsoft/DUCK/backend/ducklib/config"
	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/handlers/dictionaries"
	"github.com/Microsoft/DUCK/backend/ducklib/handlers/documents"
	"github.com/Microsoft/DUCK/backend/ducklib/handlers/rulebases"
	"github.com/Microsoft/DUCK/backend/ducklib/handlers/users"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//GetServer returns Echo instance with predefined routes
func GetServer(conf config.Configuration) *echo.Echo {
	datab, err := db.NewDatabase(*conf.DBConfig)
	if err != nil {
		switch t := err.(type) {
		case structs.HTTPError:
			e := err
			if t.Cause != nil {
				log.Printf("Database error: " + err.Error())
				e = t.Cause
			}
			panic(e)
		default:
			panic(err)
		}
	}

	JWT := conf.JwtKey
	rbd := conf.RulebaseDir

	log.Printf("Rulebase directory: " + rbd)

	checker, err := carneades.MakeComplianceCheckerPlugin(rbd)
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

	uh := users.Handler{Db: datab, JWT: JWT}
	e.POST("/login", uh.Login)
	//create sub-router for api functions
	api := e.Group("/v1")

	////User resources
	jwtMiddleware := middleware.JWT(JWT)

	users := api.Group("/users") //base URI

	//create a new user - JWT must not be required since during registration (when the user account is created) the user is not authenticated
	users.POST("", uh.PostUser)
	users.DELETE("/:id", uh.DeleteUser, jwtMiddleware) //delete a user
	users.PUT("/", uh.PutUser, jwtMiddleware)          //update a user

	dih := dictionaries.Handler{Db: datab}
	users.GET("/:id/dictionary", dih.GetUserDict, jwtMiddleware)             //get a users dictonary
	users.PUT("/:id/dictionary", dih.PutUserDict, jwtMiddleware)             //update a users dictonary
	users.GET("/:id/dictionary/:code", dih.GetDictItem, jwtMiddleware)       //get a dictonary entry
	users.PUT("/:id/dictionary/:code", dih.PutDictItem, jwtMiddleware)       //update a dictonary entry
	users.DELETE("/:id/dictionary/:code", dih.DeleteDictItem, jwtMiddleware) //delete a dictonary entry

	//data use statement document resources
	doh := documents.Handler{Db: datab}
	documents := api.Group("/documents", jwtMiddleware)    //base URI
	documents.POST("", doh.PostDoc)                        //create document
	documents.PUT("", doh.PutDoc)                          //update document
	documents.DELETE("/:docid", doh.DeleteDoc)             //delete document
	documents.GET("/:userid/summary", doh.GetDocSummaries) //return document summaries for the author
	documents.GET("/:docid", doh.GetDoc)                   //return document
	documents.POST("/copy/:docid", doh.CopyStatements)     //copies the statements from an existing Document to a new one

	//rulebase resources
	ruh := rulebases.Handler{Db: datab, WebDir: conf.WebDir, Checker: checker}
	rulebases := api.Group("/rulebases", jwtMiddleware)             //base URI
	rulebases.GET("", ruh.GetRulebases)                             //Returns a dictionary with all available Rulebases
	rulebases.PUT("/:baseid/documents", ruh.CheckDoc)               //process provided document against rulebase
	rulebases.PUT("/:baseid/documents/:documentid", ruh.CheckDocID) //process document against rulebase

	// serves the static files
	wbd := conf.WebDir

	log.Printf("Web directory: " + wbd)
	e.Static("/", wbd)

	log.Println("Server started")
	return e

}
