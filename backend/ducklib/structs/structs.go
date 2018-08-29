// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package structs

import (
	"bytes"
	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type DBConf struct {
	Location string `json:"location"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty`
}

type User struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	Password         string     `json:"password"`
	Firstname        string     `json:"firstname"`
	Lastname         string     `json:"lastname"`
	Locale           string     `json:"locale"`
	AssumptionSet    string     `json:"assumptionSet"`
	Revision         string     `json:"revision"`
	GlobalDictionary Dictionary `json:"globalDictionary"`
	//Documents []string `json:"documents"`
}

//Document struct contains among other things a set of statements, a Unique ID, a name and an owner
//The Owner field is a foreign key to a User.ID
type Document struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Revision      string      `json:"revision"`
	Owner         string      `json:"owner"`
	Locale        string      `json:"locale"`
	Description   string      `json:"description"`
	AssumptionSet string      `json:"assumptionSet"`
	Statements    []Statement `json:"statements"`
	Dictionary    Dictionary  `json:"dictionary"`
}

//UserIsOwner checks if the user ID from the JWT in the context object is the same as the user ID in the Owner field of this document
func (d *Document) IsUserOwner(c echo.Context) error {

	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return NewHTTPError("Could not access jwt", 401)
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return NewHTTPError("Could not convert jwt", 401)
	}
	id, ok := claims["id"].(string)
	if !ok {
		return NewHTTPError("Could not access user ID from JWT", 401)
	}
	if id == d.Owner {
		return nil
	}
	return NewHTTPError("User ID is not Owner ID", 401)
}

//A Statement struct represents one Statement in a document
type Statement struct {
	UseScopeCode     string           `json:"useScopeCode"`
	QualifierCode    string           `json:"qualifierCode"`
	DataCategoryCode string           `json:"dataCategoryCode"`
	DataCategories   []DataCategories `json:"dataCategories"`
	SourceScopeCode  string           `json:"sourceScopeCode"`
	ActionCode       string           `json:"actionCode"`
	ResultScopeCode  string           `json:"resultScopeCode"`
	TrackingID       string           `json:"trackingId"`
	Tag              *string          `json:"tag,omitempty"`
	Passive          bool             `json:"passive"`
}

type Operator int

const (
	AND Operator = iota
	EXCEPT
)

var operatorId = map[Operator]string{
	AND:    "and",
	EXCEPT: "except",
}

var operatorName = map[string]Operator{
	"and":    AND,
	"except": EXCEPT,
}

func (o Operator) String() string {
	return operatorId[o]
}

func (o *Operator) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(operatorId[*o])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (o *Operator) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*o = operatorName[s]
	return nil
}

type DataCategories struct {
	Op               Operator `json:"operator"`
	QualifierCode    string   `json:"qualifierCode"`
	DataCategoryCode string   `json:"dataCategoryCode"`
}

type DictionaryEntry struct {
	Value          string `json:"value"`
	//Case_1         string `json:"case_1"`
	//Case_2         string `json:"case_2"`
	Type           string `json:"type"`
	Code           string `json:"code"`
	Location       string `json:"location"`
	Category       string `json:"category"`
	DictionaryType string `json:"dictionaryType"`
}

type Dictionary map[string]DictionaryEntry

//Response represents a JSON response from the ducklib server
type Response struct {
	Ok        bool        `json:"ok"`
	Reason    *string     `json:"reason,omitempty"`
	ID        *string     `json:"id,omitempty"`
	Documents *[]Document `json:"documents,omitempty"`
}

type ComplianceResponse struct {
	Compliant   string      `json:"compliant"`
	Explanation interface{} `json:"explanation"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Rulebase struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Revision string `json:"_rev"`
}

//
type Taxonomy map[string][]struct {
	Value    string `json:"value"`
	Code     string `json:"code"`
	Category string `json:"category"`
	Fixed    bool   `json:"fixed"`
}

// HTTPError is an error with an http statuscode, it can also wrap another underlying error
type HTTPError struct {
	Err    string
	Status int
	Cause  error
}

// NewHTTPError returns a httpError which implements the Error interface and has the additional field Status for a http status code.
func NewHTTPError(err string, code int) HTTPError {
	return HTTPError{err, code, nil}
}

// WrapErrWith wraps herr around err by setting err as Cause of herr
func WrapErrWith(err error, herr HTTPError) HTTPError {
	herr.Cause = err
	return herr
}

func (e HTTPError) Error() string {
	return e.Err

}
