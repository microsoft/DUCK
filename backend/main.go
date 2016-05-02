package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type response struct {
	Message string `json:"message"`
}

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Message: "Hello World"})
}
func messageHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Message: c.Param("message")})
}

func main() {

	var webDir string

	flag.StringVar(&webDir, "webdir", "frontend", "The root directory for serving web content")
	flag.Parse()

	fmt.Println("Web root: " + webDir)

	//New echo instance
	e := echo.New()

	//set used Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//set routes for api
	e.GET("/api", helloHandler)
	e.GET("/api/:message", messageHandler)

	// serves the static files
	e.Static("/", webDir)

	//start server
	e.Run(standard.New(":3000"))

}
