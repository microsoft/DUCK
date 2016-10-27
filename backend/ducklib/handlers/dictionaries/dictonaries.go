package dictionaries

import (
	"log"
	"net/http"

	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	"github.com/labstack/echo"
)

//Handler ...
type Handler struct {
	Db *db.Database
}

//GetUserDict returns the dictionary struct on the user. if it is null, an empty one will be created and returned
//
//Context-Parameter
//	id		the id of the user whose dictionary should be returned
func (h *Handler) GetUserDict(c echo.Context) error {
	dict, err := h.Db.GetUserDict(c.Param("id"))
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	if dict == nil {
		dict = make(structs.Dictionary)
	}

	return c.JSON(http.StatusOK, dict)
}

//GetDictItem returns the dictionary entry from the users ditionary
//an error will be retuned if the dictionary does not contain the specified key
//
//Context-Parameter
//	id		the id of the user whose dictionary should be accessed
//	code	the key for the dictionary entry
func (h *Handler) GetDictItem(c echo.Context) error {
	dict, err := h.Db.GetUserDict(c.Param("id"))
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	if entry, prs := dict[c.Param("code")]; prs {
		return c.JSON(http.StatusOK, entry)
	}
	e := "Code not found in dictionary"

	return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
}

//DeleteDictItem deletes an entry from an users dict if the specified key exists
//
//Context-Parameter
//	id		the id of the user whose dictionary should be accessed
//	code	the key for the dictionary entry
//
//returns okay if the entry is not in the ditcionary anymore or never was
func (h *Handler) DeleteDictItem(c echo.Context) error {
	id := c.Param("id")
	dict, err := h.Db.GetUserDict(id)
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	delete(dict, c.Param("code"))

	err = h.Db.PutUserDict(dict, id)
	if err != nil {
		log.Printf("Error in getUserDictHandler: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusConflict, structs.Response{Ok: false, Reason: &e})
		}
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

//PutDictItem places an dictionary entry into the users dictionary
//if the key already exists the entry will be overwritten
//
//Context-Parameter
//	id				the id of the user whose dictionary should be accessed
//	code			the key for the dictionary entry
// 	in RequestBody	the DictionaryEntry
//
//returns the code if successful
func (h *Handler) PutDictItem(c echo.Context) error {
	d := new(structs.DictionaryEntry)
	if err := c.Bind(d); err != nil {
		log.Printf("Error in putDictItemHandler while trying to bind new dictionary entry to struct: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	code := c.Param("code")
	id := c.Param("id")

	dict, err := h.Db.GetUserDict(id)
	if err != nil {
		log.Printf("Error in putDictItemHandler while trying to  user dictionary from database: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	if dict == nil {
		dict = make(structs.Dictionary)
	}
	dict[code] = *d

	err = h.Db.PutUserDict(dict, id)
	if err != nil {
		log.Printf("Error in putDictItemHandler while trying to update user dictionary in database: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	return c.JSON(http.StatusOK, code)
}

//PutUserDict updates the users dictionary with a new one
//Context-Parameter
//	id				the id of the user whose dictionary should be updated
// 	in RequestBody	the new dictionary
//
//returns the new dictionary if successful
func (h *Handler) PutUserDict(c echo.Context) error {

	d := new(structs.Dictionary)
	if err := c.Bind(d); err != nil {
		log.Printf("Error in putUserDictHandler while trying to bind new dictionary to struct: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}
	id := c.Param("id")

	err := h.Db.PutUserDict(*d, id)
	if err != nil {
		log.Printf("Error in putUserDictHandler while trying to update user dictionary in database: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	nd, err := h.Db.GetUserDict(id)
	if err != nil {
		log.Printf("Error in putUserDictHandler while trying to get updated user dictionary: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	return c.JSON(http.StatusOK, nd)
}
