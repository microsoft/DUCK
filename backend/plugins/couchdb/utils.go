// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package couchdb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/Microsoft/DUCK/backend/ducklib/structs"
)

// getMap returns a map[string]interface{} containing the unmarshaled JSON from the io.Reader
func getMap(resp io.Reader) (map[string]interface{}, error) {

	content, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 502))
	}

	//TODO: remove this when not used anymore
	//fmt.Println(string(content))

	var data map[string]interface{}

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, structs.WrapErrWith(err, structs.NewHTTPError(err.Error(), 502))
	}
	//e := structs.NewHttpError(err, 404)

	return data, nil
}

func getRows(jsonbody map[string]interface{}) ([]interface{}, error) {

	rows, prs := jsonbody["rows"].([]interface{})
	//fmt.Println(jsonbody)
	if !prs || len(rows) < 1 {
		return nil, structs.NewHTTPError("No Data returned", 502)
	}
	return rows, nil
}

//FromValueMap fills the fields of a document struct with values that
//are Unmarshalled from JSON into a map
func docFromValueMap(mp map[string]interface{}) structs.Document {

	var d structs.Document
	if id, ok := mp["_id"]; ok {
		d.ID = id.(string)
	}
	if rev, ok := mp["_rev"]; ok {
		d.Revision = rev.(string)
	}
	if name, ok := mp["name"]; ok {
		d.Name = name.(string)
	}
	if owner, ok := mp["owner"]; ok {
		d.Owner = owner.(string)
	}
	if locale, ok := mp["locale"]; ok {
		d.Locale = locale.(string)
	}
	if assumptionSet, ok := mp["assumptionSet"]; ok {
		d.AssumptionSet = assumptionSet.(string)
	}
	if description, ok := mp["description"]; ok {
		d.Description = description.(string)
	}
	//Add Statements
	d.Statements = make([]structs.Statement, 0)
	if stmts, prs := mp["statements"].([]interface{}); prs {

		for _, stmt := range stmts {

			if stmt != nil {
				s := stmtFromInterfaceMap(stmt.(map[string]interface{}))
				d.Statements = append(d.Statements, s)
			}

		}
	}
	fmt.Println(d.Statements)
	//add Dictionary
	if dict, prs := mp["dictionary"].(map[string]interface{}); prs {
		d.Dictionary = dictFromInterfaceMap(dict)
	}
	return d
}

//
func stmtFromInterfaceMap(mp map[string]interface{}) structs.Statement {

	var s structs.Statement
	s.UseScopeCode = getFieldValue(mp, "useScopeCode")
	s.QualifierCode = getFieldValue(mp, "qualifierCode")
	s.DataCategoryCode = getFieldValue(mp, "dataCategoryCode")
	s.SourceScopeCode = getFieldValue(mp, "sourceScopeCode")
	s.ActionCode = getFieldValue(mp, "actionCode")
	s.ResultScopeCode = getFieldValue(mp, "resultScopeCode")
	s.TrackingID = getFieldValue(mp, "trackingId")

	s.DataCategories = make([]structs.DataCategories, 0)
	fmt.Println("DATACATEGORIES")
	if dcs, prs := mp["dataCategories"].([]interface{}); prs {
		fmt.Println(dcs)

		for _, dcmap := range dcs {
			fmt.Println(dcmap)
			var dc structs.DataCategories

			dc.DataCategoryCode = getFieldValue(dcmap.(map[string]interface{}), "dataCategoryCode")
			dc.QualifierCode = getFieldValue(dcmap.(map[string]interface{}), "qualifierCode")

			if interf, ok := dcmap.(map[string]interface{})["operator"]; ok {
				fmt.Println("Value:")
				value, ok := interf.(float64)
				fmt.Println(value)
				if ok {

					//i, err := strconv.Atoi(value)
					//if err != nil {

					dc.Op = structs.Operator(int(value))
					//}
				}
				
				fmt.Println("Valueend")
			}
			fmt.Println(dcmap.(map[string]interface{})["operator"])
			fmt.Println(dc)
			//s := stmtFromInterfaceMap(dc.(map[string]interface{}))
			s.DataCategories = append(s.DataCategories, dc)

		}
	}

	//set Tag only if it is not empty,
	//when we set tag to an empty string we cannot return null
	if tag := getFieldValue(mp, "tag"); tag != "" {
		s.Tag = &tag
	}

	s.Passive = getFieldBooleanValue(mp, "passive")

	return s
}

func getFieldValue(mp map[string]interface{}, field string) string {

	if interf, ok := mp[field]; ok {
		if value, ok := interf.(string); ok {
			return value
		}
	}
	return ""
}

func getFieldBooleanValue(mp map[string]interface{}, field string) bool {

	if interf, ok := mp[field]; ok {
		if str, ok := interf.(string); ok {
			b, err := strconv.ParseBool(str)
			if err == nil {
				return b
			}
		}
	}
	return false
}

func userFromValueMap(mp map[string]interface{}) structs.User {

	var u structs.User
	if id, ok := mp["_id"]; ok {
		u.ID = id.(string)
	}
	if rev, ok := mp["_rev"]; ok {
		u.Revision = rev.(string)
	}
	if name, ok := mp["firstname"]; ok {
		u.Firstname = name.(string)
	}
	if owner, ok := mp["lastname"]; ok {
		u.Lastname = owner.(string)
	}
	if owner, ok := mp["password"]; ok {
		u.Password = owner.(string)
	}
	if owner, ok := mp["email"]; ok {
		u.Email = owner.(string)
	}
	if locale, ok := mp["locale"]; ok {
		u.Locale = locale.(string)
	}
	if assumptionSet, ok := mp["assumptionSet"]; ok {
		u.AssumptionSet = assumptionSet.(string)
	}

	if dict, prs := mp["dictionary"].(map[string]interface{}); prs {
		u.GlobalDictionary = dictFromInterfaceMap(dict)
	}
	return u
}

func dictFromInterfaceMap(mp map[string]interface{}) structs.Dictionary {

	//Map looks like
	//dictionary:
	//	map[
	//		microsoft_excel:map[code:microsoft_excel category:1 value:Microsoft Excel type:scope]
	//		microsoft_word:map[category:1 value:Microsoft Word type:scope code:microsoft_word]]

	d := make(structs.Dictionary)

	for key, value := range mp {
		var de structs.DictionaryEntry

		value := value.(map[string]interface{})
		if code, ok := value["code"]; ok {
			de.Code = code.(string)
		}
		if tpe, ok := value["type"]; ok {
			de.Type = tpe.(string)
		}
		if val, ok := value["value"]; ok {
			de.Value = val.(string)
		}
		if category, ok := value["category"]; ok {
			de.Category = category.(string)
		}
		if dictionaryType, ok := value["dictionaryType"]; ok {
			de.DictionaryType = dictionaryType.(string)
		}

		d[key] = de
	}
	return d
}
