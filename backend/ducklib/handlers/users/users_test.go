package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Microsoft/DUCK/backend/ducklib/config"
	"github.com/Microsoft/DUCK/backend/ducklib/db"
	"github.com/Microsoft/DUCK/backend/ducklib/structs"
	_ "github.com/Microsoft/DUCK/backend/plugins/mockdb"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var (
	conf structs.Configuration
	e    *echo.Echo

	users map[string]struct {
		Pass bool         `json:"pass"`
		User structs.User `json:"user"`
	}

	//eg userIDs["user_a"]="a structs.User.ID"
	userIDs = make(map[string]string)

	uh    Handler
	JWT   []byte
	datab *db.Database
)

func TestUserHandler(t *testing.T) {

	conf = config.NewConfiguration(filepath.Join(os.Getenv("GOPATH"), "/src/github.com/Microsoft/DUCK/backend/configuration.json"))
	e = echo.New()

	JWT = []byte(conf.JwtKey)

	uh = Handler{}
	dab, err := db.NewDatabase(*conf.DBConfig)
	if err != nil {
		t.Skip("User Handler test failed; was not able to datab.Init()")
	}
	uh.Db = dab

	dat, err := ioutil.ReadFile("testdata/user.json")

	if err = json.Unmarshal(dat, &users); err != nil {
		t.Error("Testfixture User not correctly loading")
		t.Skip("No testfixtures no usertests")
	}
	//t.Logf("User %+v\n", users)
	//t.Error("AHHHHHH")
	t.Run("PostUser=1", testPostUser)
	t.Run("PostUser=2", testPostUserAgain)
	t.Run("Login=1", testLogin)
	t.Run("login=2", testWrongLogin)
	t.Run("PutUser=1", testPutUser)
	t.Run("DeleteUser=1", testDeleteUser)

}

//userhandlertests
func testPostUser(t *testing.T) {

	/*e := GetServer(conf, goPath)
	if e == nil {
		t.Fatal("Get Server Failed")
	}*/
	for key, value := range users {

		userJSON, err := json.Marshal(value.User)
		if err != nil {
			t.Errorf("Test with %s: user Post into json Marshalling not functioning", key)

		}

		req, err := http.NewRequest(echo.POST, "/users", bytes.NewReader(userJSON))
		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		err = uh.PostUser(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if value.Pass {

			if rec.Code != http.StatusCreated {
				t.Errorf("Test with %s: user creation does not return HTTP code %d but %d.", key, http.StatusCreated, rec.Code)
			} else {

				// compare with user fields since some fields are unique

				var res structs.User

				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Errorf("Test with %s: user Creation does not return valid User struct", key)
				}

				if res.Email != value.User.Email {
					t.Errorf("Test with %s: User creation returns User email %s, wants %s", key, res.Email, value.User.Email)
				}
				if res.Firstname != value.User.Firstname {
					t.Errorf("Test with %s: User creation returns User Firstname %s, wants %s", key, res.Firstname, value.User.Firstname)
				}
				if res.Lastname != value.User.Lastname {
					t.Errorf("Test with %s: User creation returns User Lastname %s, wants %s", key, res.Lastname, value.User.Lastname)
				}
				//Passwords are not retuned anymore
				/*if res.Password != value.User.Password {
					t.Errorf("Test with %s: User creation returns User Password %s, wants %s", key, res.Password, value.User.Password)
				}*/
				userIDs[key] = res.ID
				value.User.ID = res.ID
			}

		} else { //  test missing values

			if rec.Code < 400 {
				t.Errorf("Test with %s: user creation does not return a HTTP error code (>=400) but %d.", key, rec.Code)
			}
		}

	}

}

func testPostUserAgain(t *testing.T) {

	//test if already existing user is not saved again
	for key, value := range users {
		userJSON, err := json.Marshal(value.User)
		if err != nil {
			t.Errorf("Test with %s: second user Post  into json Marshalling not functioning", key)

		}
		req, err := http.NewRequest(echo.POST, "/users", bytes.NewReader(userJSON))
		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		err = uh.PostUser(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if rec.Code < 400 {
			t.Errorf("Test with %s: user creation does not return a HTTP error code (>=400) but %d.", key, rec.Code)
		}

	}

}

func testLogin(t *testing.T) {
	for key, value := range users {

		value.User.ID = userIDs[key]

		login := structs.Login{Email: value.User.Email, Password: value.User.Password}

		userJSON, err := json.Marshal(login)
		if err != nil {
			t.Errorf("Test with %s: user login json Marshal not functioning", key)

		}

		req, err := http.NewRequest(echo.POST, "/login", bytes.NewReader(userJSON))

		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		err = uh.Login(c)

		if value.Pass {
			if err != nil {
				t.Errorf("Test with %s: Error logging in: %s", key, err)
			}
			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: user Login does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
			} else {

				var dat map[string]interface{}

				//log.Println(rec.Body.String())

				if err := json.Unmarshal(rec.Body.Bytes(), &dat); err != nil {
					t.Errorf("Test with %s: user login does not return valid JSON", key)
				}

				if _, prs := dat["token"]; !prs {
					t.Errorf("Test with %s: user login does not return token", key)

				}

				if s, prs := dat["firstName"]; prs {
					if value.User.Firstname != s.(string) {
						t.Errorf("Test with %s: User login returns User Firstname %s, wants %s", key, s.(string), value.User.Firstname)
					}
				} else {
					t.Errorf("Test with %s: User login does not return Firstname", key)

				}

				if s, prs := dat["lastName"]; prs {
					if value.User.Lastname != s.(string) {
						t.Errorf("Test with %s: User login returns User Lastname %s, wants %s", key, s.(string), value.User.Lastname)
					}
				} else {
					t.Errorf("Test with %s: User login does not return Lastname", key)

				}

				if s, prs := dat["id"]; prs {
					if value.User.ID != s.(string) {
						t.Errorf("Test with %s: User login returns User id %s, wants %s", key, s.(string), value.User.ID)
					}
				} else {
					t.Errorf("Test with %s: User login does not return id", key)

				}

			}
		} else {
			if rec.Code != echo.ErrUnauthorized.Code {
				t.Errorf("Test with %s: user login does not return HTTP code %d but %d.", key, echo.ErrUnauthorized.Code, rec.Code)
			}
		}

	}
}

func testWrongLogin(t *testing.T) {
	key := "user_a"
	value := users[key]
	value.User.Password += "WrongPassword"
	logins := make([]interface{}, 4)

	logins[0] = structs.Login{Email: value.User.Email, Password: value.User.Password}
	logins[1] = structs.Login{Email: value.User.Email}
	docs := make([]structs.Document, 2)
	logins[2] = structs.Response{Documents: &docs}
	logins[3] = "teststring"

	for i, login := range logins {
		userJSON, err := json.Marshal(login)
		if err != nil {
			t.Errorf("Test with login %d: user login with wrong password json Marshal not functioning", i)

		}

		req, err := http.NewRequest(echo.POST, "/login", bytes.NewReader(userJSON))

		if err != nil {
			t.Errorf("Test with login %d: Error login with wrong password: %s", i, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		err = uh.Login(c)

		if rec.Code != echo.ErrUnauthorized.Code && rec.Code != echo.ErrNotFound.Code {
			t.Errorf("Test with login %d: wrong user login test does not return HTTP code %d or %d but %d.", i, echo.ErrUnauthorized.Code, echo.ErrNotFound.Code, rec.Code)
		}
	}
}

func testPutUser(t *testing.T) {
	for key, value := range users {
		if !value.Pass {
			continue
		}
		value.User.ID = userIDs[key]

		value.User.Firstname = fmt.Sprintf("xx%s~", value.User.Firstname)

		userJSON, err := json.Marshal(value.User)
		if err != nil {
			t.Errorf("Test with %s: user update json Marshal not functioning", key)

		}

		req, err := http.NewRequest(echo.PUT, "/users/:id", bytes.NewReader(userJSON))

		if err != nil {
			t.Errorf("Test with %s: Error creating User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
		c.SetParamNames("id")
		c.SetParamValues(value.User.ID)
		err = uh.PutUser(c)
		if err != nil {
			t.Errorf("Test with %s: Error creating User during post:%s", key, err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Test with %s: user update does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
		} else {

			var res structs.User

			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Errorf("Test with %s: user update does not return valid User struct", key)
			}

			if res.Email != value.User.Email {
				t.Errorf("Test with %s: User update returns User email %s, wants %s", key, res.Email, value.User.Email)
			}
			if res.Firstname != value.User.Firstname {
				t.Errorf("Test with %s: User update returns User Firstname %s, wants %s", key, res.Firstname, value.User.Firstname)
			}
			if res.Lastname != value.User.Lastname {
				t.Errorf("Test with %s: User update returns User Lastname %s, wants %s", key, res.Lastname, value.User.Lastname)
			}
			if res.Password != value.User.Password {
				t.Errorf("Test with %s: User update returns User Password %s, wants %s", key, res.Password, value.User.Password)
			}
		}

	}
}

func testDeleteUser(t *testing.T) {
	for key, value := range users {

		req, err := http.NewRequest(echo.DELETE, "/users/:id", nil)
		if err != nil {
			t.Errorf("Test with %s: Error deleting User: %s", key, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))

		c.SetParamNames("id")
		c.SetParamValues(userIDs[key])

		err = uh.DeleteUser(c)
		if err != nil {
			t.Errorf("Test with %s: Error deleting User during post:%s", key, err)
		}

		if value.Pass {

			if rec.Code != http.StatusOK {
				t.Errorf("Test with %s: user deletion does not return HTTP code %d but %d.", key, http.StatusOK, rec.Code)
			}

		} else { //  test missing user

			if rec.Code != http.StatusNotFound {
				t.Errorf("Test with %s: user deletion does not return HTTP code %d but %d.", key, http.StatusNotFound, rec.Code)
			}
		}

	}
}
