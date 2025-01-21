package storage_test

import (
	"election/storage"
	"election/types"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestDevStorageNewLASFSElection(t *testing.T) {
	t.Run(
		"DevStorage newLASFSElection 01",
		func(t *testing.T) {
			re1 := regexp.MustCompile(`^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2}):(?P<MM>\d{2}):(?P<SS>\d{2})Z Test1$`)

			s := storage.NewDevStorage()

			n1 := time.Now().UTC()
			id1 := s.NewLASFSElection("Test1", n1)

			if !re1.MatchString(id1) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH:MM:SSZ Test1' format: '%v'", id1)
			}
			expectedId1 := n1.Format(time.RFC3339) + " Test1"
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

			re2 := regexp.MustCompile(`^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2}):(?P<MM>\d{2}):(?P<SS>\d{2})Z Test2$`)
			n2 := time.Now().UTC()
			id2 := s.NewLASFSElection("Test2", n2)

			if !re2.MatchString(id2) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH:MM:SSZ Test2' format: '%v'", id2)
			}
			expectedId2 := n2.Format(time.RFC3339) + " Test2"
			if id1 != expectedId1 {
				t.Errorf("id '%v' does not match expected value '%v'", id2, expectedId2)
			}

			electionIDs2 := s.GetLASFSElectionsIDs()
			expectedElectionIDs2 := []string{
				id1, id2,
			}

			if !reflect.DeepEqual(electionIDs2, expectedElectionIDs2) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionIDs2, expectedElectionIDs2)
			}

			re3 := regexp.MustCompile(`^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2}):(?P<MM>\d{2}):(?P<SS>\d{2})Z Test3$`)
			n3 := time.Now().UTC()
			id3 := s.NewLASFSElection("Test3", n3)

			if !re3.MatchString(id3) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH:MM:SSZ Test3' format: '%v'", id3)
			}
			expectedId3 := n3.Format(time.RFC3339) + " Test3"
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
			m := types.LASFSMember{
				Name:     "Test",
				Password: "encryptedhash",
			}
			id, err := s.CreateLASFSMember(m)

			expectedId := "1"
			if id != expectedId {
				t.Errorf("CreateLASFSMember id '%v' does not match expected value '%v'", id, expectedId)
			}
			if err != nil {
				t.Errorf("CreateLASFSMember had unexpected error creating member '%v'", m)
			}
		},
	)
}
