package ducklib

import (
	"fmt"

	"net/http"
	"os"
	"path/filepath"
	"time"

	"log"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var goPath = os.Getenv("GOPATH")
var testData = filepath.Join(goPath, "/src/github.com/Microsoft/DUCK/testdata.json")

//route Handlers

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

func getDocSummaries(c echo.Context) error {

	docs, err := datab.GetDocumentSummariesForUser(c.Param("userid"))

	if err != nil {
		log.Printf("Error in getDocSummaries: %s", err)
		log.Println(err)
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, docs)
}

func testdataHandler(c echo.Context) error {

	if err := FillTestdata(testData); err != nil {
		log.Printf("Error in testdataHandler while trying to fill the database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusConflict, structs.Response{Ok: false, Reason: &e})

	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

/*
Document handlers
*/
func getDocHandler(c echo.Context) error {
	doc, err := datab.GetDocument(c.Param("docid"))
	if err != nil {
		log.Printf("Error in getDocHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	fmt.Printf("GET revision: %s\n", doc.Revision)
	return c.JSON(http.StatusOK, doc)
}

func copyDocHandler(c echo.Context) error {
	doc, err := datab.GetDocument(c.Param("docid"))
	if err != nil {
		log.Printf("Error in copyDocHandler trying to get old document from database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	newDoc := new(structs.Document)
	if err := c.Bind(newDoc); err != nil {
		e := err.Error()

		log.Printf("Error in copyDocHandler trying to bind newDoc: %s", err)
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	newDoc.Statements = doc.Statements

	id, err := datab.PostDocument(*newDoc)
	if err != nil {
		log.Printf("Error in copyDocHandler trying to post newDoc to database: %s", err)

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	returnDoc, err := datab.GetDocument(id)
	if err != nil {
		log.Printf("Error in copyDocHandler, trying to get newDoc: %s", err)

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, returnDoc)
}

func deleteDocHandler(c echo.Context) error {
	err := datab.DeleteDocument(c.Param("docid"))
	if err != nil {
		e := err.Error()
		log.Printf("Error in deleteDocHandler: %s", err)

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

func putDocHandler(c echo.Context) error {
	/*
		resp, err := ioutil.ReadAll(c.Request().Body())
		if err != nil {
			e := err.Error()
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
		fmt.Println(string(resp))
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false})
	*/
	doc := new(structs.Document)
	if err := c.Bind(doc); err != nil {
		e := err.Error()

		log.Printf("Error in putDocHandler while trying to bind new doc to struct: %s", err)

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	fmt.Printf("PUT revision: %s\n", doc.Revision)
	err := datab.PutDocument(*doc)
	if err != nil {
		e := err.Error()

		log.Printf("Error in putDocHandler while trying to update document in database: %s", err)

		if e == "Document update conflict." {
			return c.JSON(http.StatusConflict, structs.Response{Ok: false, Reason: &e})
		}
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	docu, err := datab.GetDocument(doc.ID)
	fmt.Printf("PUT RETURN revision: %s\n", docu.Revision) // should be the same one we once got through the document GET
	if err != nil {
		e := err.Error()

		log.Printf("Error in putDocHandler while trying to get updated document: %s", err)

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, docu)
}
func postDocHandler(c echo.Context) error {

	doc := new(structs.Document)
	if err := c.Bind(doc); err != nil {
		log.Printf("Error in postDocHandler while trying to bind new doc to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusBadRequest, structs.Response{Ok: false, Reason: &e})
	}

	id, err := datab.PostDocument(*doc)
	if err != nil {
		log.Printf("Error in postDocHandler while trying to create document in database: %s", err)

		e := err.Error()
		return c.JSON(http.StatusBadRequest, structs.Response{Ok: false, Reason: &e})
	}
	docu, err := datab.GetDocument(id)
	if err != nil {
		log.Printf("Error in postDocHandler while trying to get new document: %s", err)

		e := err.Error()
		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusCreated, docu)
}

/*
User handlers
*/

func deleteUserHandler(c echo.Context) error {
	err := datab.DeleteUser(c.Param("id"))
	if err != nil {
		log.Printf("Error in deleteUserHandler: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

func putUserHandler(c echo.Context) error {

	u := new(structs.User)
	if err := c.Bind(u); err != nil {
		log.Printf("Error in putUserHandler while trying to bind new user to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	id := c.Param("id")
	u.ID = id
	err := datab.PutUser(*u)
	if err != nil {
		log.Printf("Error in putUserHandler while trying to update user in database: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	us, err := datab.GetUser(id)
	if err != nil {
		log.Printf("Error in putUserHandler while trying to get updated user: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, us)
}
func postUserHandler(c echo.Context) error {
	newUser := new(structs.User)
	if err := c.Bind(newUser); err != nil {
		log.Printf("Error in postUserHandler while trying to bind new user to struct: %s", err)

		e := err.Error()

		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}

	id, err := datab.PostUser(*newUser)
	if err != nil {
		log.Printf("Error in postUserHandler while trying to create user in database: %s", err)

		e := err.Error()

		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}
	var u, err2 = datab.GetUser(id)
	if err2 != nil {
		log.Printf("Error in postUserHandler while trying to get new user: %s", err)

		e := err2.Error()

		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusCreated, u)
}

/*
Rulebase handlers
*/

func checkDocHandler(c echo.Context) error {
	/*
		resp, err := ioutil.ReadAll(c.Request().Body())
		if err != nil {
			e := err.Error()
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
		fmt.Println(string(resp))
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false})
	*/
	id := c.Param("baseid")

	doc := new(structs.Document)
	if err := c.Bind(doc); err != nil {
		log.Printf("Error in checkDocHandler while trying to bind document to struct: %s", err)

		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	ok, err := checker.IsCompliant(id, doc)

	//log.Printf("DOCS: %+v", docs)
	if err != nil {
		log.Printf("Error in checkDocHandler while checking for compliance: %s", err)

		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	if ok {

		return c.JSON(http.StatusOK, structs.ComplianceResponse{Ok: ok, Compliant: "COMPLIANT"})
	}

	return c.JSON(http.StatusOK, structs.ComplianceResponse{Ok: ok, Compliant: "NON_COMPLIANT"})
}

func checkDocIDHandler(c echo.Context) error {
	id := c.Param("baseid")
	docid := c.Param("documentid")

	doc, err := datab.GetDocument(docid)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to get document from database: %s", err)

		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	ok, err := checker.IsCompliant(id, &doc)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while checking for compliance: %s", err)
		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	fmt.Printf("Compliant: %t", ok)
	if ok {
		return c.JSON(http.StatusOK, structs.ComplianceResponse{Ok: ok, Compliant: "COMPLIANT"})
	}

	return c.JSON(http.StatusOK, structs.ComplianceResponse{Ok: ok, Compliant: "NON_COMPLIANT"})

}
func getRulebasesHandler(c echo.Context) error {
	//log.Printf("Rulebases: %+v", checker.RuleBases)
	if len(checker.RuleBases) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, checker.RuleBases)
}

// loginHandler handles the login Process
func loginHandler(c echo.Context) error {
	/*	resp, err := ioutil.ReadAll(c.Request().Body())
		if err != nil {
			e := err.Error()
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
		fmt.Println("USER: ")
		fmt.Println(string(resp))
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false})*/
	u := new(structs.Login)
	if err := c.Bind(u); err != nil {
		log.Printf("Error in loginHandler trying to bind user to struct: %s", err)
		return err
	}

	id, pw, err := datab.GetLogin(u.Email) //TODO compare with encrypted pw
	if err != nil {
		log.Printf("Error in loginHandler trying to get login info for userMail %s: %s", u.Email, err)

		e := err.Error()

		return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &e})
	}
	//log.Printf("id: %s, pw: %s", id, pw)
	if u.Password == pw {

		user, err := datab.GetUser(id)
		if err != nil {
			log.Printf("Error in loginHandler trying to get user info: %s", err)

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
			log.Printf("Error in loginHandler: %s", err)
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
	reason := "Passwords do not match"
	log.Printf("Error in loginHandler: %s", reason)

	return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &reason})

}
