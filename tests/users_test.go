package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	// Build a user and then fetch by its ID
	var ub testutils.UserBuilder
	u, err := ub.New().WithName("Temp").Create()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+testutils.Itoa(u.ID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		User models.JsonUser `json:"user"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, u.ID, resp.User.ID)
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


func TestUsersPosts_Show_Existing(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	// Build a user and then fetch by its ID
	var ub testutils.UserBuilder
	u, err := ub.New().WithName("Temp").Create()
	assert.NoError(t, err)

	// Build posts for the user
	var pb testutils.PostBuilder
	for i := range 5 {
		_, err = pb.New().WithTitle("Temp"+strconv.Itoa(i)).WithBody("Temp"+strconv.Itoa(i)).WithUserID(u.ID).Create()
		assert.NoError(t, err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+testutils.Itoa(u.ID)+"/posts", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		User models.JsonUser `json:"user"`
		Posts []models.JsonPost `json:"posts"`
	}
	json_body := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, json_body)
	assert.Equal(t, u.ID, resp.User.ID)
	assert.Equal(t, 5, len(resp.Posts))
}


func TestUsersPosts_User_Not_Found(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+testutils.Itoa(999999)+"/posts", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

}