package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
curl -H "Content-Type: application/json" -d '{"user_id":"12345", "venue_id": "1"}' http://localhost:8080/create_user_favorite
*/
func UserFavoriteCreate(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var userFav schema.UserFavorite
		err := decoder.Decode(&userFav)
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err)
			return
		}
		defer r.Body.Close()

		id, err := db.CreateUserFavorite(userFav)
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
curl -H "Content-Type: application/json" -d '{"user_id":"1234"}' http://localhost:8080/user_favorites
*/
func UserFavoritesList(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		keys, ok := r.URL.Query()["user_id"]

		if !ok || len(keys) < 1 {
			log.Println("Url Param 'user_id' is missing")
			return
		}
		uid := keys[0]
		u := schema.UserFavorite{UserID: uid}

		favorites, err := db.UserFavoritesList(u)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Data []schema.Venue `json:"data"`
		}
		writeJSON(w, http.StatusCreated, envelope{favorites})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" http://localhost:8080/user_favorites/:venue_id/:user_id
*/
func UserFavoritesGet(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		vid, err := strconv.Atoi(vars["venue_id"])
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		uid := vars["user_id"]
		u := schema.UserFavorite{UserID: uid, VenueID: vid}

		favorite, err := db.UserFavoritesGet(u)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Data schema.Venue `json:"data"`
		}
		writeJSON(w, http.StatusCreated, envelope{favorite})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" http://localhost:8080/user_favorites/:id
*/
func UserFavoritesRemove(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		err = db.UserFavoritesDelete(id)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		writeJSON(w, http.StatusAccepted, nil)
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"user_id":"1234", "venue_id": "1"}' http://localhost:8080/remove_user_favorite
*/
// func UserFavoriteDelete(db data.Database) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		decoder := json.NewDecoder(r.Body)
// 		var userFav schema.UserFavorite
// 		err := decoder.Decode(&userFav)
// 		if err != nil {
// 			writeError(w, http.StatusUnprocessableEntity, err)
// 			return
// 		}
// 		defer r.Body.Close()

// 		id, err := db.DeleteUserFavorite(userFav)
// 		if err != nil {
// 			writeError(w, http.StatusUnprocessableEntity, err)
// 			return
// 		}

// 		type envelope struct {
// 			Status string `json:"status"`
// 			Result int    `json:"result"`
// 		}
// 		writeJSON(w, http.StatusCreated, envelope{http.StatusText(http.StatusCreated), id})
// 	})
// }

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
curl -H "Content-Type: application/json" -d '{"name":"Panzano", "address": "909 17th St", "city": "Denver", "zip": "80202", "state": "CO", "image": "http://coloradobites.com/wp-content/uploads/2015/05/panzanococktail1.jpg", "country": "USA"}' http://localhost:8080/create_venue
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

		menu := schema.Menu{VenueID: id}
		menuID, err := db.CreateMenu(menu)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Status  string `json:"status"`
			VenueID int    `json:"venue_id"`
			MenuID  int    `json:"menu_id"`
		}
		writeJSON(w, http.StatusCreated, envelope{http.StatusText(http.StatusCreated), id, menuID})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"Popular"}' http://localhost:8080/create_venue_list
*/
func VenueListCreate(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var venueList schema.VenueList
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		if err := r.Body.Close(); err != nil {
			return
		}
		if err := json.Unmarshal(body, &venueList); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		id, err := db.CreateVenueList(venueList)
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

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"name":"Popular"}' http://localhost:8080/venue_list
*/
func VenueListGet(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		keys, ok := r.URL.Query()["name"]

		if !ok || len(keys) < 1 {
			log.Println("Url Param 'name' is missing")
			return
		}
		name := keys[0]
		venueList := schema.VenueList{Name: name}

		vl, err := db.VenueListGet(venueList)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		venues, err := db.VenuesByList(vl.ID)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Data []schema.Venue `json:"data"`
		}
		writeJSON(w, http.StatusCreated, envelope{venues})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json"  http://localhost:8080/venue/1
*/
func VenueGet(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		v := schema.Venue{ID: id}
		venue, err := db.VenueGet(v)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Data schema.Venue `json:"data"`
		}
		writeJSON(w, http.StatusCreated, envelope{venue})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json"  http://localhost:8080/menu_items?venue_id=1
*/
func MenuItemsGet(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		keys, ok := r.URL.Query()["venue_id"]

		if !ok || len(keys) < 1 {
			log.Println("Url Param 'venue_id' is missing")
			return
		}
		id, err := strconv.Atoi(keys[0])
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		m := schema.Menu{VenueID: id}
		menus, err := db.MenuItemsGet(m)
		if err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}

		type envelope struct {
			Data []schema.MenuItem `json:"data"`
		}
		writeJSON(w, http.StatusCreated, envelope{menus})
	})
}

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"venue_name":"Panzano", "venue_list_name": "Popular"}' http://localhost:8080/venue_list_add
*/
func VenueListAdd(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var vl schema.VenueListAdd
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		if err := r.Body.Close(); err != nil {
			return
		}

		if err := json.Unmarshal(body, &vl); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		id, err := db.VenueListAdd(vl)
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

/*
Test with this curl command:
curl -H "Content-Type: application/json" -d '{"menu_id": 1, "category": "Drink", "price": 5.00, "description": "LOCAL DRAFT BEERS"}' http://localhost:8080/add_menu_item
*/
func MenuItemAdd(db data.Database) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var m schema.MenuItem
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			fmt.Println("test: ", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		if err := r.Body.Close(); err != nil {
			return
		}

		if err := json.Unmarshal(body, &m); err != nil {
			fmt.Println("test2: ", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		id, err := db.AddToMenu(m)
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
