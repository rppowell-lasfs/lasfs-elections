package api_test

import (
	"election/cmd/api"
	"election/storage"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetElectionNoNominees(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n := time.Now().UTC()
	id := s.NewLASFSElection("Test1", n, nil).ID

	req, _ := http.NewRequest("GET", fmt.Sprintf("/election/%s", id), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedData := fmt.Sprintf(`{"election:":{"election":"%s","position":"Test1","status":"","nominees":[]}}`, id)
	assert.Equal(t, expectedData, string(responseData))
}

func TestGetElectionWithNominees(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n := time.Now().UTC()
	id := s.NewLASFSElection("Test1", n, []string{"Alpha", "Bravo", "Charlie"}).ID

	req, _ := http.NewRequest("GET", fmt.Sprintf("/election/%s", id), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedData := fmt.Sprintf(`{"election:":{"election":"%s","position":"Test1","status":"","nominees":["Alpha","Bravo","Charlie"]}}`, id)
	assert.Equal(t, expectedData, string(responseData))
}
