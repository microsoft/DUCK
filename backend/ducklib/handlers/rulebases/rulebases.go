// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package rulebases

import (
	"log"
	"net/http"

	"github.com/Microsoft/DUCK/backend/ducklib/carneades"
	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/labstack/echo"
)

//Handler ...
type Handler struct {
	Db      *db.Database
	WebDir  string
	Checker *carneades.ComplianceCheckerPlugin
}

//CheckDoc checks the document against a rulebase for compliance
//
//Context-Parameter
//	baseid			the id of the rulebase
// 	in RequestBody	the document
func (h *Handler) CheckDoc(c echo.Context) error {
	id := c.Param("baseid")
	doc := new(structs.Document)
	if err := c.Bind(doc); err != nil {
		log.Printf("Error in checkDocHandler while trying to bind document to struct: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	normalizer, err := carneades.NewNormalizer(*doc, h.Db, h.WebDir)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to normalize document : %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	normDoc, err := normalizer.GetNormalized()
	if err != nil {
		log.Printf("Error in checkDocHandler while normalizing: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	ok, exp, err := h.Checker.IsCompliant(id, normDoc)

	if err != nil {
		log.Printf("Error in checkDocHandler while checking for compliance: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	flatExp := carneades.FoldExplanation(exp)

	if ok {
		return c.JSON(http.StatusOK, structs.ComplianceResponse{Compliant: "COMPLIANT", Explanation: flatExp})
	}
	return c.JSON(http.StatusOK, structs.ComplianceResponse{Compliant: "NON_COMPLIANT", Explanation: flatExp})
}

//CheckDocID checks a document from the database against a rulebase for compliance
//
//Context-Parameter
//	baseid		the id of the rulebase
// 	documentid	the id of the document
func (h *Handler) CheckDocID(c echo.Context) error {
	id := c.Param("baseid")
	docid := c.Param("documentid")

	doc, err := h.Db.GetDocument(docid)

	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to get document from database: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	normalizer, err := carneades.NewNormalizer(doc, h.Db, h.WebDir)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while trying to normalize document : %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	normDoc, err := normalizer.GetNormalized()
	if err != nil {
		log.Printf("Error in checkDocIDHandler while normalizing: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	ok, exp, err := h.Checker.IsCompliant(id, normDoc)
	if err != nil {
		log.Printf("Error in checkDocIDHandler while checking for compliance: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	if ok {
		return c.JSON(http.StatusOK, structs.ComplianceResponse{Compliant: "COMPLIANT", Explanation: exp})
	}

	return c.JSON(http.StatusOK, structs.ComplianceResponse{Compliant: "NON_COMPLIANT", Explanation: exp})

}

//GetRulebases returns a list of  all loaded rulebases
func (h *Handler) GetRulebases(c echo.Context) error {
	//if we have no loaded rulebases return Error
	if len(h.Checker.RuleBases) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, h.Checker.RuleBases)
}
