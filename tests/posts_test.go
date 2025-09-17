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

func TestPosts_Index_OK(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Posts []models.JsonPost `json:"posts"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.GreaterOrEqual(t, len(resp.Posts), 2)
}

func TestPosts_Show_Existing(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	// Build a post and then fetch by its ID
	p, err := testutils.NewPostBuilder().WithTitle("Seeded").WithBody("From test").Create()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/"+testutils.Itoa(p.ID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Post models.JsonPost `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, p.ID, resp.Post.ID)
}

func TestPosts_Show_NotFound(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/99999", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPosts_Create_Valid(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"title":"New Title","body":"New Body"}`)
	req, _ := http.NewRequest("POST", "/posts/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Post models.JsonPost `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "New Title", resp.Post.Title)
}

func TestPosts_Create_BadRequest(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"title":"","body":""}`)
	req, _ := http.NewRequest("POST", "/posts/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPosts_Update_Existing(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	// Create to get a known ID
	p, err := testutils.NewPostBuilder().WithTitle("Temp").WithBody("Temp").Create()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	body := []byte(`{"title":"Updated","body":"Updated Body"}`)
	req, _ := http.NewRequest("PATCH", "/posts/"+testutils.Itoa(p.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Post models.JsonPost `json:"post"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Updated", resp.Post.Title)
}

func TestPosts_Update_NotFound(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	body := []byte(`{"title":"Updated","body":"Updated Body"}`)
	req, _ := http.NewRequest("PATCH", "/posts/99999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPosts_Delete_Existing(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	// Create to get a known ID
	p, err := testutils.NewPostBuilder().WithTitle("Temp").WithBody("Temp").Create()
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/"+testutils.Itoa(p.ID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify it is gone
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/posts/"+testutils.Itoa(p.ID), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPosts_Delete_NotFound(t *testing.T) {
	cleanup, err := testutils.BeginTxWithSeeds()
	assert.NoError(t, err)
	defer cleanup()

	router := NewRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/99999", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
