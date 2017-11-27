package route

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kernkw/hhapp/internal/data"
	"github.com/kernkw/hhapp/internal/data/datamock"
	"github.com/kernkw/hhapp/internal/schema"
)

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}

func TestUserCreate(t *testing.T) {
	wantID := 1234567
	mockStore := &datamock.Mock{
		CreateUser_: func(user schema.User) (int, error) {
			return wantID, nil
		},
	}

	u := schema.User{
		UserName: "test",
		Password: "password",
		Email:    "test@domain.com",
	}

	jsonU, err := json.Marshal(u)
	checkError(err, t)
	req, err := http.NewRequest("POST", "/create_account", bytes.NewReader(jsonU))
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserCreate(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := fmt.Sprintf(`{"status":"Created","result":%d}`, wantID)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUserCreate_bad_input(t *testing.T) {
	mockStore := &datamock.Mock{}

	u := schema.User{}

	jsonU, err := json.Marshal(u)
	checkError(err, t)
	req, err := http.NewRequest("POST", "/create_account", bytes.NewReader(jsonU))

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserCreate(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

}

func TestUserCreate_bad_error(t *testing.T) {
	wantErr := data.ErrDuplicateEntry
	mockStore := &datamock.Mock{
		CreateUser_: func(user schema.User) (int, error) {
			return 0, wantErr
		},
	}

	u := schema.User{
		UserName: "test",
		Password: "password",
		Email:    "test@domain.com",
	}

	jsonU, err := json.Marshal(u)
	checkError(err, t)
	req, err := http.NewRequest("POST", "/create_account", bytes.NewReader(jsonU))

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserCreate(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	expected := fmt.Sprintf(`{"status":"%v"}`, wantErr)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected error: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUserFavoriteCreate(t *testing.T) {
	wantID := 1234567
	mockStore := &datamock.Mock{
		CreateUserFavorite_: func(userFav schema.UserFavorite) (int, error) {
			return wantID, nil
		},
	}

	u := schema.UserFavorite{
		UserID:  "12345",
		VenueID: 1,
	}

	jsonU, err := json.Marshal(u)
	checkError(err, t)
	req, err := http.NewRequest("POST", "/create_user_favorite", bytes.NewReader(jsonU))
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserFavoriteCreate(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := fmt.Sprintf(`{"status":"Created","result":%d}`, wantID)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUserFavoriteCreate_error(t *testing.T) {
	wantErr := errors.New("Some error")
	mockStore := &datamock.Mock{
		CreateUserFavorite_: func(userFav schema.UserFavorite) (int, error) {
			return 0, wantErr
		},
	}

	u := schema.UserFavorite{
		UserID:  "12345",
		VenueID: 1,
	}

	jsonU, err := json.Marshal(u)
	checkError(err, t)
	req, err := http.NewRequest("POST", "/create_user_favorite", bytes.NewReader(jsonU))
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserFavoriteCreate(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	expected := fmt.Sprintf(`{"status":"%s"}`, wantErr)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUserFavoritesList(t *testing.T) {
	wants := []schema.Venue{
		schema.Venue{
			ID:      1,
			Name:    "test",
			Address: "12345 test street",
			City:    "Noname",
			State:   "CO",
			Zip:     "123456",
			Country: "USA",
			Image:   "http://someimage",
		},
		schema.Venue{
			ID:      2,
			Name:    "test2",
			Address: "12345 test street",
			City:    "Noname",
			State:   "CO",
			Zip:     "123456",
			Country: "USA",
			Image:   "http://someimage",
		},
	}
	mockStore := &datamock.Mock{
		UserFavoritesList_: func(userFav schema.UserFavorite) ([]schema.Venue, error) {
			return wants, nil
		},
	}

	uid := "12345"
	req, err := http.NewRequest("GET", fmt.Sprintf("/user_favorites?user_id=%s", uid), nil)
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserFavoritesList(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	wantedJSONResponse, err := json.Marshal(wants)
	checkError(err, t)

	expected := fmt.Sprintf(`{"data":%v}`, string(wantedJSONResponse))
	fmt.Println(rr.Body)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUserFavoritesList_error(t *testing.T) {
	wantErr := errors.New("Some error")
	mockStore := &datamock.Mock{
		UserFavoritesList_: func(userFav schema.UserFavorite) ([]schema.Venue, error) {
			return nil, wantErr
		},
	}

	uid := "12345"
	req, err := http.NewRequest("GET", fmt.Sprintf("/user_favorites?user_id=%s", uid), nil)
	checkError(err, t)

	rr := httptest.NewRecorder()

	http.HandlerFunc(UserFavoritesList(mockStore)).
		ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	expected := fmt.Sprintf(`{"status":"%s"}`, wantErr)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
