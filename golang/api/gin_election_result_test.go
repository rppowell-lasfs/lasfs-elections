package api_test

import (
	"election/api"
	"election/storage"
	"election/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetElectionResult1Vote(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n1 := time.Now().UTC()
	e1 := s.NewLASFSElection("Test", n1, []string{"Alpha", "Bravo", "Charlie"})
	id1 := e1.ID

	s.CreateLASFSMember("TestVoter", "TestPassword")

	b1 := types.NewLASFSBallot("TestVoter", "TestVoter", []string{"Alpha", "Bravo", "Charlie"})
	s.AddLASFSBallot(id1, b1)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/electionresults/%s", id1), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")
	t.Log(w.Body.String())

	expectedStruct := struct {
		ElectionID       string         `json:"election"`
		ElectionPosition string         `json:"position"`
		Nominees         []string       `json:"nominees"`
		BallotCount      int            `json:"ballotcount"`
		Ballots          map[string]int `json:"votes"`
	}{
		ElectionID:       id1,
		ElectionPosition: "Test",
		Nominees:         []string{"Alpha", "Bravo", "Charlie"},
		BallotCount:      1,
		Ballots: map[string]int{
			"Alpha":       1,
			"deadballots": 0,
		},
	}
	expectedJSON, _ := json.Marshal(expectedStruct)
	assert.Equal(t, string(expectedJSON), w.Body.String())
}

func TestGetElectionResult3VoteS(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n1 := time.Now().UTC()
	e1 := s.NewLASFSElection("Test", n1, []string{"Alpha", "Bravo", "Charlie"})
	id1 := e1.ID

	s.CreateLASFSMember("TestVoter1", "TestPassword")
	b1 := types.NewLASFSBallot("TestVoter1", "TestVoter1", []string{"Alpha", "Bravo", "Charlie"})
	s.AddLASFSBallot(id1, b1)

	s.CreateLASFSMember("TestVoter2", "TestPassword")
	b2 := types.NewLASFSBallot("TestVoter2", "TestVoter2", []string{"Alpha", "Bravo", "Charlie"})
	s.AddLASFSBallot(id1, b2)

	s.CreateLASFSMember("TestVoter3", "TestPassword")
	b3 := types.NewLASFSBallot("TestVoter3", "TestVoter3", []string{"Alpha", "Bravo", "Charlie"})
	s.AddLASFSBallot(id1, b3)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/electionresults/%s", id1), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")
	t.Log(w.Body.String())

	expectedStruct := struct {
		ElectionID       string         `json:"election"`
		ElectionPosition string         `json:"position"`
		Nominees         []string       `json:"nominees"`
		BallotCount      int            `json:"ballotcount"`
		Ballots          map[string]int `json:"votes"`
	}{
		ElectionID:       id1,
		ElectionPosition: "Test",
		Nominees:         []string{"Alpha", "Bravo", "Charlie"},
		BallotCount:      3,
		Ballots: map[string]int{
			"Alpha":       3,
			"deadballots": 0,
		},
	}
	expectedJSON, _ := json.Marshal(expectedStruct)
	assert.Equal(t, string(expectedJSON), w.Body.String())
}
