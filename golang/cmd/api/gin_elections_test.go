package api_test

import (
	"election/cmd/api"
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

	n1 := time.Now().UTC()
	id1 := s.NewLASFSElection("Test1", n1)

	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	r.GET("/elections", ginHandler.GetLASFSElections)

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

	n1 := time.Now().UTC()
	id1 := s.NewLASFSElection("Test1", n1)
	n2 := time.Now().UTC()
	id2 := s.NewLASFSElection("Test2", n2)
	n3 := time.Now().UTC()
	id3 := s.NewLASFSElection("Test3", n3)

	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	r.GET("/elections", ginHandler.GetLASFSElections)

	req, _ := http.NewRequest("GET", "/elections", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedData := `{"elections":["` + id1 + `","` + id2 + `","` + id3 + `"]}`
	assert.Equal(t, expectedData, string(responseData))
}
