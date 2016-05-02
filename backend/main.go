package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "duck" && password == "duck" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		token.Claims["name"] = "Jon Snow"
		token.Claims["admin"] = true
		token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
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

	//username: duck  password:duck
	e.POST("/login", loginHandler)
	//create sub-router for api functions
	api := e.Group("/api")

	//set routes for api
	api.GET("", helloHandler)

	api.GET("/:message", messageHandler)

	//create restricted sub-router
	restricted := api.Group("/restricted", middleware.JWTAuth([]byte("secret")))
	//set restricted routes for api
	restricted.GET("", helloHandler)
	restricted.GET("/:message", messageHandler)

	// serves the static files
	e.Static("/", webDir)

	//start server
	e.Run(standard.New(":3000"))

}
