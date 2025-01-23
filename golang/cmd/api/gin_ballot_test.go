package api_test

import (
	"bytes"
	"election/cmd/api"
	"election/storage"
	"election/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// GetBallotForElectionByMember
func TestGetBallotForElectionByMember(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n := time.Now().UTC()
	e := s.NewLASFSElection("Test1", n, []string{"Alpha", "Bravo", "Charlie"})
	id := e.ID

	member, _ := s.CreateLASFSMember("TestVoter", "TestPassword")

	b := types.NewLASFSBallot("TestVoter", "TestVoter", []string{"Alpha", "Bravo", "Charlie"})
	s.AddLASFSBallot(id, b)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/vote/%s/%s", id, member.ID), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedData := fmt.Sprintf(`{"election":"%s","id":"TestVoter","nominees":["Alpha","Bravo","Charlie"]}`, id)
	assert.Equal(t, expectedData, string(responseData))
}

// GetBallotForElectionByMember
func TestPostBallotForElectionByMember1(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n := time.Now().UTC()
	e := s.NewLASFSElection("Test1", n, []string{"Alpha", "Bravo", "Charlie"})
	eid := e.ID

	member, _ := s.CreateLASFSMember("TestVoter", "encryptedhash")

	// b := types.LASFSBallot{
	// 	VoterName: "TestVoter",
	// 	Nominees: []types.NomineeEntry{
	// 		{NomineeName: "Alpha", IsValid: true},
	// 		{NomineeName: "Bravo", IsValid: true},
	// 		{NomineeName: "Charlie", IsValid: true},
	// 	},
	// 	IsDead: false,
	// }
	// s.AddLASFSBallotToLASFSElection(id, b)

	jsonValue, _ := json.Marshal([]string{"Alpha", "Bravo", "Charlie"})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/vote/%s/%s", eid, member.ID), bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code, "StatusAccepted")

	fmt.Printf("%s", w.Body.String())
	// w := httptest.NewRecorder()
	// r.ServeHTTP(w, req)

	// assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	// responseData, _ := io.ReadAll(w.Body)
	// expectedData := fmt.Sprintf(`{"election:":{"election":"%s","position":"Test1","status":"","nominees":["Alpha","Bravo","Charlie"]}}`, id)
	// assert.Equal(t, expectedData, string(responseData))
}
