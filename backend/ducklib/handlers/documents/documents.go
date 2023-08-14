// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package documents

import (
	"log"
	"net/http"

	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//Handler ...
type Handler struct {
	Db *db.Database
}

//GetDocSummaries returns ID and name for each Document that has the field owner with a specified userID
//
//Context-Parameter:
//	userid		a userid string which is showing to the user that owns the documents
func (h *Handler) GetDocSummaries(c echo.Context) error {

	//check if user in JWT & userid-param is same
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		e := "Could not access jwt"
		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		e := "Could not convert jwt"
		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}
	id, ok := claims["id"].(string)
	if !ok {
		e := "Could not access user ID from JWT"
		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}
	if id != c.Param("userid") {
		e := "User ID is not Param ID"
		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}

	docs, err := h.Db.GetDocumentSummariesForUser(c.Param("userid"))

	if err != nil {
		log.Printf("Error in getDocSummaries: %s", err)
		log.Println(err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, docs)
}

//GetDoc resturns a document if it exists in the database
//
//Context-Parameter:
//	docid		a docid string which is pointing to the wanted document
func (h *Handler) GetDoc(c echo.Context) error {
	doc, err := h.Db.GetDocument(c.Param("docid"))
	if err != nil {
		log.Printf("Error in getDocHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	err = doc.IsUserOwner(c)
	if err != nil {
		log.Printf("Error in getDocHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}
	return c.JSON(http.StatusOK, doc)
}

//CopyStatements  copies the Statements from one document
//in the database to a new one
//
//Context-Parameter:
//	docid			a docid string which is pointing to the wanted document that is to be copied from
//
//	in RequestBody		containing a new document without an statements
//
//Returns the new document if successful
func (h *Handler) CopyStatements(c echo.Context) error {
	doc, err := h.Db.GetDocument(c.Param("docid"))
	if err != nil {
		log.Printf("Error in copyStatementsHandler trying to get old document from database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	err = doc.IsUserOwner(c)
	if err != nil {
		log.Printf("Error in getDocHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}

	newDoc := new(structs.Document)
	if err := c.Bind(newDoc); err != nil {
		e := err.Error()

		log.Printf("Error in copyStatementsHandler trying to bind newDoc: %s", err)
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	newDoc.Statements = doc.Statements

	id, err := h.Db.PostDocument(*newDoc)
	if err != nil {
		log.Printf("Error in copyStatementsHandler trying to post newDoc to database: %s", err)

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	returnDoc, err := h.Db.GetDocument(id)
	if err != nil {
		log.Printf("Error in copyStatementsHandler, trying to get newDoc: %s", err)

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, returnDoc)
}

//DeleteDoc deletes a document if it exists in the database
//
//Context-Parameter:
//	docid		a docid string which is pointing to the wanted document
func (h *Handler) DeleteDoc(c echo.Context) error {

	doc, err := h.Db.GetDocument(c.Param("docid"))
	if err != nil {
		log.Printf("Error in deleteDocHandler trying to get document from database: %s", err)
		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	err = doc.IsUserOwner(c)
	if err != nil {
		log.Printf("Error in deleteDocHandler: %s", err)
		e := err.Error()

		return c.JSON(http.StatusForbidden, structs.Response{Ok: false, Reason: &e})
	}

	err = h.Db.DeleteDocument(c.Param("docid"))
	if err != nil {
		e := err.Error()
		log.Printf("Error in deleteDocHandler: %s", err)

		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

//PutDoc replaces a document in the database with a newer version if both have the same revision number
//
//Context-Parameter
//	in RequestBody		the new version of the document
//
//Returns the new version if successful
func (h *Handler) PutDoc(c echo.Context) error {
	doc := new(structs.Document)
	if err := c.Bind(doc); err != nil {
		e := err.Error()
		log.Printf("Error in putDocHandler while trying to bind new doc to struct: %s", err)
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	err := h.Db.PutDocument(*doc)
	if err != nil {
		e := err.Error()
		log.Printf("Error in putDocHandler while trying to update document in database: %s", err)
		if e == "Document update conflict." {
			return c.JSON(http.StatusConflict, structs.Response{Ok: false, Reason: &e})
		}
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	docu, err := h.Db.GetDocument(doc.ID)

	if err != nil {
		e := err.Error()
		log.Printf("Error in putDocHandler while trying to get updated document: %s", err)
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}
	return c.JSON(http.StatusOK, docu)
}

//PostDoc creates a new structs.Document entry in the database
//
//Context-Parameter
//	in RequestBody:		the new Document
//
//Returns the new Document if successful
func (h *Handler) PostDoc(c echo.Context) error {

	doc := new(structs.Document)
	if err := c.Bind(doc); err != nil {
		log.Printf("Error in postDocHandler while trying to bind new doc to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusBadRequest, structs.Response{Ok: false, Reason: &e})
	}

	id, err := h.Db.PostDocument(*doc)
	if err != nil {
		log.Printf("Error in postDocHandler while trying to create document in database: %s", err)

		e := err.Error()
		return c.JSON(http.StatusBadRequest, structs.Response{Ok: false, Reason: &e})
	}
	docu, err := h.Db.GetDocument(id)
	if err != nil {
		log.Printf("Error in postDocHandler while trying to get new document: %s", err)

		e := err.Error()
		return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusCreated, docu)

}
