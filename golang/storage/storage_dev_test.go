package storage_test

import (
	"election/storage"
	"election/types"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDevStorageNewLASFSElection(t *testing.T) {
	t.Run(
		"DevStorage newLASFSElection 01",
		func(t *testing.T) {

			regString := `^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2})-(?P<MM>\d{2})-(?P<SS>\d{2})\.(?P<nnnnnnnnn>\d{9})Z`

			re1 := regexp.MustCompile(regString + `_Test1$`)

			s := storage.NewDevStorage()

			n1 := time.Now().UTC()
			e1 := s.NewLASFSElection("Test1", n1, nil)
			id1 := e1.ID

			if !re1.MatchString(id1) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH-MM-SS.nnnnnnnnnZ_Test1' format: '%v'", id1)
			}
			expectedId1 := n1.Format(`2006-01-02T15-04-05.000000000Z0700`) + "_Test1"
			if id1 != expectedId1 {
				t.Errorf("id '%v' does not match expected value '%v'", id1, expectedId1)
			}

			electionIDs1 := s.GetLASFSElectionsIDs()
			expectedElectionIDs1 := []string{
				id1,
			}

			if !reflect.DeepEqual(electionIDs1, expectedElectionIDs1) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionIDs1, expectedElectionIDs1)
			}

			re2 := regexp.MustCompile(regString + `_Test_2$`)
			n2 := time.Now().UTC()
			e2 := s.NewLASFSElection("Test 2", n2, nil)
			id2 := e2.ID

			if !re2.MatchString(id2) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH-MM-SS.nnnnnnnnnZ_Test_2' format: '%v'", id2)
			}
			expectedId2 := n2.Format(`2006-01-02T15-04-05.000000000Z0700`) + "_Test_2"
			if id2 != expectedId2 {
				t.Errorf("id '%v' does not match expected value '%v'", id2, expectedId2)
			}

			electionIDs2 := s.GetLASFSElectionsIDs()
			expectedElectionIDs2 := []string{
				id1, id2,
			}

			if !reflect.DeepEqual(electionIDs2, expectedElectionIDs2) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionIDs2, expectedElectionIDs2)
			}

			re3 := regexp.MustCompile(regString + `_Test_3$`)
			n3 := time.Now().UTC()
			e3 := s.NewLASFSElection("Test:3", n3, nil)
			id3 := e3.ID

			if !re3.MatchString(id3) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH-MM-SS.nnnnnnnnnZ_Test_3' format: '%v'", id3)
			}
			expectedId3 := n3.Format(`2006-01-02T15-04-05.000000000Z0700`) + "_Test_3"
			if id3 != expectedId3 {
				t.Errorf("id '%v' does not match expected value '%v'", id3, expectedId3)
			}

			electionIDs3 := s.GetLASFSElectionsIDs()
			expectedElectionIDs3 := []string{
				id1, id2, id3,
			}

			if !reflect.DeepEqual(electionIDs3, expectedElectionIDs3) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionIDs3, expectedElectionIDs3)
			}
		},
	)
}

func TestDevStorageLASFSMember(t *testing.T) {
	t.Run(
		"DevStorage new LASFSMember 01",
		func(t *testing.T) {
			s := storage.NewDevStorage()
			var id string
			m, err := s.CreateLASFSMember("Test", "encryptedhash")
			if err != nil {
				t.Error("CreateLASFSMember had unexpected error creating member")
			} else {
				id = m.ID
			}
			expectedId := "Test"
			if id != expectedId {
				t.Errorf("CreateLASFSMember id '%v' does not match expected value '%v'", id, expectedId)
			}
		},
	)
}

func TestAddLASFSBallotToLASFSElection1(t *testing.T) {
	s := storage.NewDevStorage()
	n := time.Now().UTC()
	e := s.NewLASFSElection("Test1", n, []string{"Alpha", "Bravo", "Charlie"})
	id := e.ID

	s.CreateLASFSMember("TestVoter", "TestPassword")

	expectedDevStorage := storage.DevStorage{
		LASFSElections: map[string]*types.LASFSElection{
			id: {
				ID:              id,
				Position:        "Test1",
				Nominees:        []string{"Alpha", "Bravo", "Charlie"},
				LASFSBallots:    make([]*types.LASFSBallot, 0),
				DateTimeCreated: n,
			},
		},
		LASFSMembers: map[string]types.LASFSMember{
			"TestVoter": {ID: "TestVoter", Name: "TestVoter", Password: "TestPassword"},
		},
		LASFSMemberIDs: []string{"TestVoter"},
	}
	assert.Equal(t, &expectedDevStorage, s)
}

func TestAddLASFSBallotToLASFSElection2(t *testing.T) {
	s := storage.NewDevStorage()
	n := time.Now().UTC()
	e := s.NewLASFSElection("Test1", n, []string{"Alpha", "Bravo", "Charlie"})
	id := e.ID

	s.CreateLASFSMember("TestVoter", "TestPassword")

	b := types.LASFSBallot{
		VoterID:   "TestVoter",
		VoterName: "TestVoter",
		Nominees: []types.NomineeEntry{
			{NomineeName: "Alpha", IsValid: true},
			{NomineeName: "Bravo", IsValid: true},
			{NomineeName: "Charlie", IsValid: true},
		},
		IsDead: false,
	}

	e, err := s.AddLASFSBallot(id, b)
	assert.Equal(t, nil, err)

	expectedLASFSElection := types.LASFSElection{
		ID:       id,
		Position: "Test1",
		Nominees: []string{"Alpha", "Bravo", "Charlie"},
		LASFSBallots: []*types.LASFSBallot{
			{
				VoterID:   "TestVoter",
				VoterName: "TestVoter",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Alpha", IsValid: true},
					{NomineeName: "Bravo", IsValid: true},
					{NomineeName: "Charlie", IsValid: true},
				},
				IsDead: false,
			},
		},
		DateTimeCreated: n,
	}

	assert.Equal(t, &expectedLASFSElection, e)

	expectedDevStorage := storage.DevStorage{
		LASFSElections: map[string]*types.LASFSElection{
			id: {
				ID:       id,
				Position: "Test1",
				Nominees: []string{"Alpha", "Bravo", "Charlie"},
				LASFSBallots: []*types.LASFSBallot{
					{
						VoterID:   "TestVoter",
						VoterName: "TestVoter",
						Nominees: []types.NomineeEntry{
							{NomineeName: "Alpha", IsValid: true},
							{NomineeName: "Bravo", IsValid: true},
							{NomineeName: "Charlie", IsValid: true},
						},
						IsDead: false,
					},
				},
				DateTimeCreated: n,
			},
		},
		LASFSMembers: map[string]types.LASFSMember{
			"TestVoter": {ID: "TestVoter", Name: "TestVoter", Password: "TestPassword"},
		},
		LASFSMemberIDs: []string{"TestVoter"},
	}
	assert.Equal(t, &expectedDevStorage, s)
}
