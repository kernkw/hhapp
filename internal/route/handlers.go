package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kernkw/hhapp/internal/data"
	"github.com/kernkw/hhapp/internal/schema"
)

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"username":"test-user", "password": "password", "email": "kyle.kern@sendgrid.com"}' http://localhost:8080/create_account
*/
func UserCreate(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var user schema.User
		err := decoder.Decode(&user)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		defer r.Body.Close()

		err = user.Validate()
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		err = user.HashPassword()
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		id, err := db.CreateUser(user)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}

		type envelope struct {
			Status string `json:"status"`
			Result int    `json:"result"`
		}
		writeJSON(w, http.StatusCreated, envelope{http.StatusText(http.StatusCreated), id})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"username":"test-user", "password": "password"}' http://localhost:8080/authenticate
*/
func UserLogin(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var inuser schema.User
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := r.Body.Close(); err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := json.Unmarshal(body, &inuser); err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		dbuser, err := db.GetUser(inuser)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if !inuser.Authorized(dbuser) {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		type envelope struct {
			Status string `json:"status"`
		}
		writeJSON(w, http.StatusOK, envelope{http.StatusText(http.StatusOK)})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"Panzano", "address": "909 17th St", "city": "Denver", "zip": "80202", "state": "CO", "country": "USA"}' http://localhost:8080/create_venue
*/
func VenueCreate(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var venue schema.Venue
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		if err := r.Body.Close(); err != nil {
			return
		}
		if err := json.Unmarshal(body, &venue); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		id, err := db.CreateVenue(venue)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Status string `json:"status"`
			Result int    `json:"result"`
		}
		writeJSON(w, http.StatusCreated, envelope{http.StatusText(http.StatusCreated), id})
	})
}

func writeResponse(w http.ResponseWriter, obj map[string]string) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		w.Write(b)
	}
}

func writeError(w http.ResponseWriter, code int, err error) {
	type envelope struct {
		Status string `json:"status"`
	}
	if err == nil {
		err = errors.New(http.StatusText(code))
	}
	writeJSON(w, code, envelope{err.Error()})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(code)
	w.Write(b)
}
