package ducklib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var goPath = os.Getenv("GOPATH")

var testData = filepath.Join(goPath, "/src/github.com/Metaform/duck/testdata.json")

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

	dat, err := ioutil.ReadFile(testData)

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

/*
Document handlers
*/
func getDocHandler(c echo.Context) error {
	doc, err := datab.GetDocument(c.Param("docid"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, doc)
}
func deleteDocHandler(c echo.Context) error {
	err := datab.DeleteDocument(c.Param("docid"))
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func putDocHandler(c echo.Context) error {

	resp, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	data := Document{}
	json.Unmarshal(resp, &data)

	err = datab.PutDocument(data.ID, resp)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	doc, err := datab.GetDocument(data.ID)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, doc)
}
func postDocHandler(c echo.Context) error {

	req, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	id, err := datab.PostDocument(req)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}
	doc, err := datab.GetDocument(id)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, doc)
}

/*
User handlers
*/

func deleteUserHandler(c echo.Context) error {
	err := datab.DeleteUser(c.Param("id"))
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func putUserHandler(c echo.Context) error {

	resp, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}
	id := c.Param("id")
	err = datab.PutUser(id, resp)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	doc, err := datab.GetUser(id)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, doc)
}
func postUserHandler(c echo.Context) error {

	req, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	id, err := datab.PostUser(req)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}
	doc, err := datab.GetUser(id)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, doc)
}

/*
Ruleset handlers
*/

func deleteRsHandler(c echo.Context) error {
	err := datab.DeleteRuleset(c.Param("id"))
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, Response{Ok: true})
}

func putRsHandler(c echo.Context) error {

	resp, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}
	id := c.Param("id")
	err = datab.PutRuleset(id, resp)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	doc, err := datab.GetRuleset(id)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, doc)
}
func postRsHandler(c echo.Context) error {

	req, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	id, err := datab.PostRuleset(req)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}
	doc, err := datab.GetRuleset(id)
	if err != nil {
		e := err.Error()
		return c.JSON(http.StatusNotFound, Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, doc)
}

// loginHandler handles the login Process
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
		token.Claims["id"] = user.ID
		token.Claims["permissions"] = 1024 //FIXME
		token.Claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(JWT))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token":     t,
			"firstName": user.Firstname,
			"lastName":  user.Lastname,
			"id":        user.ID,
			"locale":    user.Locale,
		})
	}
	return echo.ErrUnauthorized

}
