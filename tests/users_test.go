package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"rest_api/models"
	"rest_api/tests/testutils"

	"github.com/stretchr/testify/assert"
)

func TestUsers_Show_Existing(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		User models.JsonUser `json:"user"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, uint(1), resp.User.ID)
}

func TestUsers_Show_NotFound(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/99999", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUsers_Create_Valid(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"name":"Charlie"}`)
	req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		User models.JsonUser `json:"user"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Charlie", resp.User.Name)
}

func TestUsers_Create_BadRequest(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"name":""}`)
	req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
