package ducklib

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//structs

var datab *Database

//route Handlers

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func getDocSummaries(c echo.Context) error {

	docs, err := datab.GetDocumentSummariesForUser(c.Param("userid"))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, docs)
}

func testdataHandler(c echo.Context) error {

	dat, err := ioutil.ReadFile("testdata.json")

	var e string
	if err != nil {
		e = err.Error()
		return c.JSON(http.StatusExpectationFailed, Response{Ok: false, Reason: &e})

	}
	if err := FillTestdata(dat); err != nil {
		e = err.Error()
		return c.JSON(http.StatusConflict, Response{Ok: false, Reason: &e})

	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func getDocHandler(c echo.Context) error {
	doc, err := datab.GetDocument(c.Param("docid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, doc)
}
func deleteHandler(c echo.Context) error {
	err := datab.Delete(c.Param("id"))
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func putHandler(c echo.Context) error {

	resp, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	err = datab.PutEntry(c.Param("id"), resp)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}
func postHandler(c echo.Context) error {

	resp, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	err = datab.PostEntry(resp)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}
func loginHandler(c echo.Context) error {

	u := new(Login)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, pw, _ := datab.GetLogin(u.Username) //TODO compare with encrypted pw

	if u.Password == pw {

		user, err := datab.GetUser(id)
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		token.Claims["firstName"] = user.Firstname
		token.Claims["lastName"] = user.Lastname
		token.Claims["id"] = user.Id
		token.Claims["permissions"] = 1024 //FIXME
		token.Claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token":     t,
			"firstName": user.Firstname,
			"lastName":  user.Lastname,
			"id":        user.Id,
		})
	}
	return echo.ErrUnauthorized

}

//GetServer returns Echo instance with predefined routes
func GetServer(webDir string, jwtKey []byte) *echo.Echo {

	datab = NewDatabase()
	datab.Init()

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

	//username: duck  password:duck
	e.POST("/login", loginHandler)
	e.GET("/loadtestdata", testdataHandler)
	//create sub-router for api functions
	api := e.Group("/v1")

	////User resources
	users := api.Group("/users", middleware.JWT(jwtKey)) //base URI

	users.POST("/", postHandler)
	users.DELETE("/:id", deleteHandler)
	users.PUT("/:id", helloHandler)

	//data use statement document resources
	documents := api.Group("/documents") //base URI
	//	documents := api.Group("/documents", middleware.JWT(jwtKey))                                          //base URI
	documents.POST("", postHandler)                    //create document
	documents.PUT("/:id", putHandler)                  //update document
	documents.DELETE("/:id", deleteHandler)            //delete document
	documents.GET("/:userid/summary", getDocSummaries) //return document summaries for the author
	documents.GET("/:userid/:docid", getDocHandler)    //return document for the author ?

	//ruleset resources
	rulesets := api.Group("/rulesets", middleware.JWT(jwtKey))  //base URI
	rulesets.POST("/", helloHandler)                            //create a ruleset
	rulesets.DELETE("/:id", deleteHandler)                      //delete a ruleset
	rulesets.PUT("/:setid", helloHandler)                       //update a ruleset
	rulesets.PUT("/:setid/documents", helloHandler)             //process provided document against ruleset
	rulesets.PUT("/:setid/documents/:documentid", helloHandler) //process document against ruleset

	// serves the static files
	e.Static("/", webDir)

	return e

}
