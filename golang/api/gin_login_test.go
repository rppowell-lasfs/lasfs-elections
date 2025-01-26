package api_test

import (
	"bytes"
	"election/api"
	"election/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run(
		"gin Tests Signup and Login",
		func(t *testing.T) {
			s := storage.NewDevStorage()
			ginHandler := api.NewGinHandler(s)

			r := gin.Default()
			api.SetupGin(ginHandler, r)

			type CreateLASFSMemberPayload struct {
				Name     string `json:"name"`
				Password string `json:"password"`
			}
			createLASFSMemberPayload := CreateLASFSMemberPayload{
				Name:     "TestName",
				Password: "TestPassword",
			}
			jsonValue, _ := json.Marshal(createLASFSMemberPayload)

			req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code, "HTTP StatusCreated")
			expectedData := `{"id":"TestName"}`
			assert.Equal(t, expectedData, w.Body.String())

			assert.Contains(t, s.LASFSMembers, createLASFSMemberPayload.Name)
		},
	)
}
