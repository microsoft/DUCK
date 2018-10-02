// Data Use Statement Compliance Checker (DUCK)
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package users

import (
	"log"
	"net/http"
	"time"

	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

//Handler ...
type Handler struct {
	Db  *db.Database
	JWT []byte
}

//DeleteUser deletes an existing user from the database
//
//Context-Parameter
//	id		the id of the user who should be deleted
func (h *Handler) DeleteUser(c echo.Context) error {
	err := h.Db.DeleteUser(c.Param("id"))
	if err != nil {
		log.Printf("Error in deleteUserHandler: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, structs.Response{Ok: true})
}

//PutUser replaces a user in the database with a newer version if both have the same revision number
//
//Context-Parameter
//	in RequestBody		the new version of the user
//
//Returns the new version if successful
func (h *Handler) PutUser(c echo.Context) error {

	u := new(structs.User)
	if err := c.Bind(u); err != nil {
		log.Printf("Error in putUserHandler while trying to bind new user to struct: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	err := h.Db.PutUser(*u)
	if err != nil {
		log.Printf("Error in putUserHandler while trying to update user in database: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	us, err := h.Db.GetUser(u.ID)
	if err != nil {
		log.Printf("Error in putUserHandler while trying to get updated user: %s", err)

		e := err.Error()
		return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
	}

	return c.JSON(http.StatusOK, us)
}

//PostUser creates a new structs.User entry in the database
//
//Context-Parameter
//	in RequestBody:		the new User
//
//Returns the new User if successful
func (h *Handler) PostUser(c echo.Context) error {

	newUser := new(structs.User)
	if err := c.Bind(newUser); err != nil {
		log.Printf("Error in postUserHandler while trying to bind new user to struct: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
		}
	}
	//no password no user
	if newUser.Password == "" {
		e := "No password submitted"
		return c.JSON(http.StatusBadRequest, structs.Response{Ok: false, Reason: &e})
	}

	//TODO: should this happen here or in db.Database.PostUser ?

	//hash password
	password := []byte(newUser.Password)
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error in postUserHandler while hashing password: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
		}
	}
	newUser.Password = string(hashedPassword)

	id, err := h.Db.PostUser(*newUser)
	if err != nil {
		log.Printf("Error in postUserHandler while trying to create user in database: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
		}
	}
	u, err := h.Db.GetUser(id)
	if err != nil {
		log.Printf("Error in postUserHandler while trying to get new user: %s", err)

		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusInternalServerError, structs.Response{Ok: false, Reason: &e})
		}
	}
	//don't show password hash to frontend
	u.Password = ""
	return c.JSON(http.StatusCreated, u)

}

//Login handles the login Process
func (h *Handler) Login(c echo.Context) error {
	u := new(structs.Login)
	if err := c.Bind(u); err != nil {
		log.Printf("Error in loginHandler trying to bind user to struct: %s", err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusNotFound, structs.Response{Ok: false, Reason: &e})
		}
	}

	id, hashedpw, err := h.Db.GetLogin(u.Email)
	if err != nil {
		log.Printf("Error in loginHandler trying to get login info for userMail %s: %s", u.Email, err)
		e := err.Error()
		switch t := err.(type) {
		case structs.HTTPError:
			return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
		default:
			return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &e})
		}
	}

	correct := true
	err = bcrypt.CompareHashAndPassword([]byte(hashedpw), []byte(u.Password))
	if err != nil {
		correct = (u.Password == hashedpw)
	}

	if correct {

		user, err := h.Db.GetUser(id)
		if err != nil {
			log.Printf("Error in loginHandler trying to get user info for userMail %s: %s", u.Email, err)
			e := err.Error()
			switch t := err.(type) {
			case structs.HTTPError:
				return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
			default:
				return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &e})
			}
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["firstName"] = user.Firstname
		claims["lastName"] = user.Lastname
		claims["id"] = user.ID
		claims["permissions"] = 1024 //FIXME
		claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(h.JWT))
		if err != nil {
			log.Printf("Error in loginHandler: %s", err)
			e := err.Error()
			switch t := err.(type) {
			case structs.HTTPError:
				return c.JSON(t.Status, structs.Response{Ok: false, Reason: &e})
			default:
				return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &e})
			}
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
	log.Printf("Error in loginHandler for userMail %s: %s", u.Email, reason)

	return c.JSON(http.StatusUnauthorized, structs.Response{Ok: false, Reason: &reason})

}
