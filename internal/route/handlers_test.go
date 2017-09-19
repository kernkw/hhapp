package route

import (
	"bytes"
	"encoding/json"
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
