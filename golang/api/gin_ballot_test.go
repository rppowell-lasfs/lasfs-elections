package api_test

import (
	"bytes"
	"election/api"
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

func TestPostBallotForElectionByMember1(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n := time.Now().UTC()
	e := s.NewLASFSElection("Test1", n, []string{"Alpha", "Bravo", "Charlie"})
	eid := e.ID

	member, _ := s.CreateLASFSMember("TestVoter", "encryptedhash")

	jsonPayload, _ := json.Marshal([]string{"Alpha", "Bravo", "Charlie"})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/vote/%s/%s", eid, member.ID), bytes.NewBuffer(jsonPayload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code, "StatusAccepted")

	t.Log(w.Body.String())

	expectedStruct := struct {
		ElectionID string   `json:"election"`
		VoterID    string   `json:"id"`
		Nominees   []string `json:"nominees"`
	}{
		ElectionID: eid,
		VoterID:    "TestVoter",
		Nominees:   []string{"Alpha", "Bravo", "Charlie"},
	}
	expectedJSON, _ := json.Marshal(expectedStruct)
	assert.Equal(t, string(expectedJSON), w.Body.String())

}

//GetBallotForElectionByMember

func TestGetLASFSElectionBallots(t *testing.T) {
	s := storage.NewDevStorage()
	ginHandler := api.NewGinHandler(s)

	r := gin.Default()
	api.SetupGin(ginHandler, r)

	n1 := time.Now().UTC()
	e1 := s.NewLASFSElection("Election1", n1, []string{"Alpha1", "Bravo1", "Charlie1"})
	eid1 := e1.ID

	member1, _ := s.CreateLASFSMember("TestVoter1", "TestPassword1")

	b1 := types.NewLASFSBallot(member1.ID, member1.Name, []string{"Alpha1", "Bravo1", "Charlie1"})
	s.AddLASFSBallot(eid1, b1)

	member2, _ := s.CreateLASFSMember("TestVoter2", "TestPassword2")

	b2 := types.NewLASFSBallot(member2.ID, member2.Name, []string{"Alpha1", "Bravo1", "Charlie1"})
	s.AddLASFSBallot(eid1, b2)

	member3, _ := s.CreateLASFSMember("TestVoter3", "TestPassword3")

	b3 := types.NewLASFSBallot(member3.ID, member3.Name, []string{"Alpha1", "Bravo1", "Charlie1"})
	s.AddLASFSBallot(eid1, b3)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/votes/%s", eid1), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "StatusOK")

	responseData, _ := io.ReadAll(w.Body)
	expectedStruct := struct {
		ElectionID       string   `json:"election"`
		ElectionPosition string   `json:"position"`
		ElectionStatus   string   `json:"status"`
		Nominees         []string `json:"nominees"`
		Ballots          []struct {
			VoterID   string   `json:"voterid"`
			VoterName string   `json:"votername"`
			Nominees  []string `json:"nominees"`
		} `json:"ballots"`
	}{
		ElectionID:       eid1,
		ElectionPosition: "Election1",
		Nominees:         []string{"Alpha1", "Bravo1", "Charlie1"},
		Ballots: []struct {
			VoterID   string   `json:"voterid"`
			VoterName string   `json:"votername"`
			Nominees  []string `json:"nominees"`
		}{
			{VoterID: member1.ID, VoterName: member1.Name, Nominees: []string{"Alpha1", "Bravo1", "Charlie1"}},
			{VoterID: member2.ID, VoterName: member2.Name, Nominees: []string{"Alpha1", "Bravo1", "Charlie1"}},
			{VoterID: member3.ID, VoterName: member3.Name, Nominees: []string{"Alpha1", "Bravo1", "Charlie1"}},
		},
	}
	expectedJSON, _ := json.Marshal(expectedStruct)
	t.Log(string(responseData))
	assert.Equal(t, string(expectedJSON), string(responseData))
}
