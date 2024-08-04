package api

import (
	"encoding/json"
	"enigma-protocol-go/pkg/db"
	"enigma-protocol-go/pkg/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
)

const DATABSE_PATH = "test.db"

func setup() *httprouter.Router {
	dbopts := &db.DatabaseOpts{
		Driver: "sqlite3",
		Uri:    DATABSE_PATH,
	}

	router, err := NewRouter(dbopts)
	if err != nil {
		panic(err)
	}

	return router
}

func cleanup() {
	os.Remove(DATABSE_PATH)
}

func TestIndexAPI(t *testing.T) {
	router := setup()
	defer cleanup()

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong Status")
	}
}

func TestNewUser(t *testing.T) {
	router := setup()
	defer cleanup()

	tests := []struct {
		name      string
		publicKey string
	}{
		{"test1", "random-public-key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/login/"+tt.publicKey, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
			}

			var res models.LoginResponse
			if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			req, _ = http.NewRequest("GET", "/connect/"+res.ID, nil)
			rr = httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
			}

			var res2 models.ConnectResponse
			if err := json.NewDecoder(rr.Body).Decode(&res2); err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if res2.Publickey != tt.publicKey {
				t.Errorf("Expected public key %v, but got %v", tt.publicKey, res2.Publickey)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	router := setup()
	defer cleanup()

	userId := "random-user"
	req, _ := http.NewRequest("GET", "/connect/"+userId, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, but got %v", http.StatusNotFound, status)
	}
}
