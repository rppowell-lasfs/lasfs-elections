package api_test

import (
	"election/api"
	"election/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetElectionsIDs1(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n1 := time.Now().UTC()
	e1 := s.NewLASFSElection("Test1", n1, nil)
	id1 := e1.ID

	req, _ := http.NewRequest("GET", "/elections", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedData := `{"elections":["` + id1 + `"]}`
	assert.Equal(t, expectedData, string(responseData))
}

func TestGetElectionsIDs3(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n1 := time.Now().UTC()
	e1 := s.NewLASFSElection("Test1", n1, nil)
	id1 := e1.ID
	n2 := time.Now().UTC()
	e2 := s.NewLASFSElection("Test 2", n2, nil)
	id2 := e2.ID
	n3 := time.Now().UTC()
	e3 := s.NewLASFSElection("Test:3", n3, nil)
	id3 := e3.ID

	req, _ := http.NewRequest("GET", "/elections", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedData := `{"elections":["` + id1 + `","` + id2 + `","` + id3 + `"]}`
	assert.Equal(t, expectedData, string(responseData))
}
