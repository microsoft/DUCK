package ducklib

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

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

//helloHandler returns just Hello world with StatusOK.
func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}

//getDocSummaries returns ID and name for each Document that has the field owner with a specified userID
//
//Context-Parameter:
//	userid		a userid string which is showing to the user that owns the documents
func getDocSummaries(c echo.Context) error {

	docs, err := datab.GetDocumentSummariesForUser(c.Param("userid"))

	if err != nil {
		log.Printf("Error in getDocSummaries: %s", err)
		log.Println(err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, docs)
}

//testdataHandler initializes the import of testing data from the file testdata.json into the database
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

//getDocHandler resturns a document if it exists in the database
//
//Context-Parameter:
//	docid		a docid string which is pointing to the wanted document
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

//copyStatementsHandler  copies the Statements from one document
//in the database to a new one
//
//Context-Parameter:
//	docid			a docid string which is pointing to the wanted document that is to be copied from
//
//	in RequestBody		containing a new document without an statements
//
//Returns the new document if successful
func copyStatementsHandler(c echo.Context) error {
	doc, err := datab.GetDocument(c.Param("docid"))
	if err != nil {
		log.Printf("Error in copyStatementsHandler trying to get old document from database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	newDoc := new(structs.Document)
	if err := c.Bind(newDoc); err != nil {
		e := err.Error()

		log.Printf("Error in copyStatementsHandler trying to bind newDoc: %s", err)
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	newDoc.Statements = doc.Statements

	id, err := datab.PostDocument(*newDoc)
	if err != nil {
		log.Printf("Error in copyStatementsHandler trying to post newDoc to database: %s", err)

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	returnDoc, err := datab.GetDocument(id)
	if err != nil {
		log.Printf("Error in copyStatementsHandler, trying to get newDoc: %s", err)

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, returnDoc)
}

//deleteDocHandler deletes a document if it exists in the database
//
//Context-Parameter:
//	docid		a docid string which is pointing to the wanted document
func deleteDocHandler(c echo.Context) error {
	err := datab.DeleteDocument(c.Param("docid"))
	if err != nil {
		e := err.Error()
		log.Printf("Error in deleteDocHandler: %s", err)

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

//putDocHandler replaces a document in the database with a newer version if both have the same revision number
//
//Context-Parameter
//	in RequestBody		the new version of the document
//
//Returns the new version if successful
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

//postDocHandler creates a new structs.Document entry in the database
//
//Context-Parameter
//	in RequestBody:		the new Document
//
//Returns the new Document if successful
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

//deleteUserHandler deletes an existing user from the database
//
//Context-Parameter
//	id		the id of the user who should be deleted
func deleteUserHandler(c echo.Context) error {
	err := datab.DeleteUser(c.Param("id"))
	if err != nil {
		log.Printf("Error in deleteUserHandler: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

//putUserHandler replaces a user in the database with a newer version if both have the same revision number
//
//Context-Parameter
//	in RequestBody		the new version of the user
//
//Returns the new version if successful
func putUserHandler(c echo.Context) error {

	u := new(structs.User)
	if err := c.Bind(u); err != nil {
		log.Printf("Error in putUserHandler while trying to bind new user to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	err := datab.PutUser(*u)
	if err != nil {
		log.Printf("Error in putUserHandler while trying to update user in database: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	us, err := datab.GetUser(u.ID)
	if err != nil {
		log.Printf("Error in putUserHandler while trying to get updated user: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, us)
}

//postUserHandler creates a new structs.User entry in the database
//
//Context-Parameter
//	in RequestBody:		the new User
//
//Returns the new User if successful
func postUserHandler(c echo.Context) error {

	newUser := new(structs.User)
	if err := c.Bind(newUser); err != nil {
		log.Printf("Error in postUserHandler while trying to bind new user to struct: %s", err)

		e := err.Error()

		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}
	//hash password
	password := []byte(newUser.Password)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error in postUserHandler while hashing password: %s", err)
		e := err.Error()
		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}
	newUser.Password = string(hashedPassword)

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
	//don't show password hash to frontend
	u.Password = ""
	return c.JSON(http.StatusCreated, u)

}

//getUserDictHandler returns the dictionary struct on the user. if it is null, an empty one will be created and returned
//
//Context-Parameter
//	id		the id of the user whose dictionary should be returned
func getUserDictHandler(c echo.Context) error {
	dict, err := datab.GetUserDict(c.Param("id"))
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	if dict == nil {
		dict = make(structs.Dictionary)
	}

	return c.JSON(http.StatusOK, dict)
}

//getUserDictHandler returns the dictionary entry from the users ditionary
//an error will be retuned if the dictionary does not contain the specified key
//
//Context-Parameter
//	id		the id of the user whose dictionary should be accessed
//	code	the key for the dictionary entry
func getDictItemHandler(c echo.Context) error {
	dict, err := datab.GetUserDict(c.Param("id"))
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	if entry, prs := dict[c.Param("code")]; prs {
		return c.JSON(http.StatusOK, entry)
	}
	e := "Code not found in dictionary"

	return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
}

//deleteDictItemHandler deletes an entry from an users dict if the specified key exists
//
//Context-Parameter
//	id		the id of the user whose dictionary should be accessed
//	code	the key for the dictionary entry
//
//returns okay if the entry is not in the ditcionary anymore or never was
func deleteDictItemHandler(c echo.Context) error {
	id := c.Param("id")
	dict, err := datab.GetUserDict(id)
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	delete(dict, c.Param("code"))

	err = datab.PutUserDict(dict, id)
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusConflict, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

//putDictItemHandler places an dictionary entry into the users dictionary
//if the key already exists the entry will be overwritten
//
//Context-Parameter
//	id				the id of the user whose dictionary should be accessed
//	code			the key for the dictionary entry
// 	in RequestBody	the DictionaryEntry
//
//returns the code if successful
func putDictItemHandler(c echo.Context) error {
	d := new(structs.DictionaryEntry)
	if err := c.Bind(d); err != nil {
		log.Printf("Error in putDictItemHandler while trying to bind new dictionary entry to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	code := c.Param("code")
	id := c.Param("id")

	dict, err := datab.GetUserDict(id)
	if err != nil {
		log.Printf("Error in putDictItemHandler while trying to  user dictionary from database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	if dict == nil {
		dict = make(structs.Dictionary)
	}
	dict[code] = *d

	err = datab.PutUserDict(dict, id)
	if err != nil {
		log.Printf("Error in putDictItemHandler while trying to update user dictionary in database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, code)
}

//putUserDictHandler updates the users dictionary with a new one
//Context-Parameter
//	id				the id of the user whose dictionary should be updated
// 	in RequestBody	the new dictionary
//
//returns the new dictionary if successful
func putUserDictHandler(c echo.Context) error {

	d := new(structs.Dictionary)
	if err := c.Bind(d); err != nil {
		log.Printf("Error in putUserDictHandler while trying to bind new dictionary to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	id := c.Param("id")

	err := datab.PutUserDict(*d, id)
	if err != nil {
		log.Printf("Error in putUserDictHandler while trying to update user dictionary in database: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	nd, err := datab.GetUserDict(id)
	if err != nil {
		log.Printf("Error in putUserDictHandler while trying to get updated user dictionary: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, nd)
}

/*
Rulebase handlers
*/

//checkDocHandler checks the document against a rulebase for compliance
//
//Context-Parameter
//	baseid			the id of the rulebase
// 	in RequestBody	the document
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
	normalizer, err := NewNormalizer(*doc, datab)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to normalize document : %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	normDoc, err := normalizer.CreateDict()
	if err != nil {
		log.Printf("Error in checkDocHandler while normalizing: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	ok, err := checker.IsCompliant(id, normDoc)
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

//checkDocIDHandler checks a document from the database against a rulebase for compliance
//
//Context-Parameter
//	baseid		the id of the rulebase
// 	documentid	the id of the document
func checkDocIDHandler(c echo.Context) error {
	id := c.Param("baseid")
	docid := c.Param("documentid")

	doc, err := datab.GetDocument(docid)

	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to get document from database: %s", err)

		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	normalizer, err := NewNormalizer(doc, datab)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to normalize document : %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	normDoc, err := normalizer.CreateDict()
	if err != nil {
		log.Printf("Error in checkDocIDHandler while normalizing: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	ok, err := checker.IsCompliant(id, normDoc)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while checking for compliance: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	if ok {
		return c.JSON(http.StatusOK, structs.ComplianceResponse{Ok: ok, Compliant: "COMPLIANT"})
	}

	return c.JSON(http.StatusOK, structs.ComplianceResponse{Ok: ok, Compliant: "NON_COMPLIANT"})

}

//getRulebasesHandler returns a list of  all loaded rulebases
func getRulebasesHandler(c echo.Context) error {
	//if we have no loaded rulebases return Error
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

	id, hashedpw, err := datab.GetLogin(u.Email) //TODO compare with encrypted pw
	if err != nil {
		log.Printf("Error in loginHandler trying to get login info for userMail %s: %s", u.Email, err)

		e := err.Error()

		return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &e})
	}
	//log.Printf("id: %s, pw: %s", id, pw)

	correct := true
	err = bcrypt.CompareHashAndPassword([]byte(hashedpw), []byte(u.Password))
	if err != nil {
		correct = (u.Password == hashedpw)
	}

	if correct {

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
