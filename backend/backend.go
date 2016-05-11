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

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type configuration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Message: "Hello World"})
}
func messageHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Message: c.Param("message")})
}

func loginHandler(c echo.Context) error {

	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}

	if u.Username == "duck" && u.Password == "duck" {

		//set vars -> we should get this from the DB
		firstName := "Duck"
		lastName := "Goose"
		id := "e6eb5f0a-2ec0-4f79-b0c9-df6e6bb66032"

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		token.Claims["firstName"] = firstName
		token.Claims["lastName"] = lastName
		token.Claims["id"] = id
		token.Claims["permissions"] = 1024 //FIXME
		token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token":     t,
			"firstName": firstName,
			"lastName":  lastName,
			"id":        id,
		})
	}

	return echo.ErrUnauthorized
}

func main() {

	var webDir string
	//jwtKey := []byte("secret")

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
	api := e.Group("/v1")
	//set routes for api
	api.GET("", helloHandler)
	api.GET("/:message", messageHandler)

	////User resources
	//users := api.Group("/users", middleware.JWT(jwtKey))
	//
	//users.POST("", helloHandler)
	//users.DELETE("/:id", helloHandler)
	//users.PUT("/:id", helloHandler)
	//
	////data use statement document resources
	//documents := api.Group("/documents", middleware.JWT(jwtKey))
	//documents.POST("", helloHandler)
	//
	////ruleset resources
	//rulesets := api.Group("/rulesets", middleware.JWT(jwtKey))
	//rulesets.POST("", helloHandler)
	//
	////create restricted sub-router
	//restricted := api.Group("/restricted", middleware.JWT(jwtKey))
	////set restricted routes for api
	//restricted.GET("", helloHandler)
	//restricted.GET("/:message", messageHandler)

	// serves the static files
	e.Static("/", webDir)

	//start server
	e.Run(standard.New(":3000"))

}
