package api_test

import (
	"bytes"
	"election/cmd/api"
	"election/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestTwo(t *testing.T) {
	t.Run(
		"gin Tests Signup and Login",
		func(t *testing.T) {
			s := storage.NewDevStorage()
			ginHandler := api.NewGinHandler(s)

			r := gin.Default()
			r.POST("/signup", ginHandler.Signup)
			r.POST("/login", ginHandler.Login)

			type CreateLASFSMemberPayload struct {
				Name     string
				Password string
			}
			createLASFSMemberPayload := CreateLASFSMemberPayload{
				Name:     "TestName",
				Password: "TestPassword",
			}
			jsonValue, _ := json.Marshal(createLASFSMemberPayload)

			req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code, "StatusCreated")

			// mockResponse := `{"name":"1"}`
			// responseData, _ := io.ReadAll(w.Body)
			// assert.Equal(t, mockResponse, responseData)
			// t.Errorf("received '%v'", w)

		},
	)
}
