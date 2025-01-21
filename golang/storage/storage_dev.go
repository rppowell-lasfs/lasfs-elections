package storage

import (
	"election/types"
	"errors"
	"sort"
	"strconv"
	"time"
)

type DevStorage struct {
	LASFSElections map[string]types.LASFSElection
	Members        []types.LASFSMember
}

func NewDevStorage() *DevStorage {
	return &DevStorage{
		LASFSElections: make(map[string]types.LASFSElection),
	}
}

func (s *DevStorage) NewLASFSElection(position string, createdDateTime time.Time) string {
	e := types.NewLASFSElection(position, createdDateTime)
	id := e.ID
	e.DateTimeCreated = createdDateTime
	s.LASFSElections[id] = *e
	return id
}

func (s *DevStorage) GetLASFSElectionsIDs() []string {
	keys := make([]string, 0, len(s.LASFSElections))

	for k := range s.LASFSElections {
		keys = append(keys, s.LASFSElections[k].ID)
	}
	sort.Strings(keys)
	return keys
}

func (s *DevStorage) GetLASFSElectionByID(id string) types.LASFSElection {
	return s.LASFSElections[id]
}

func (s *DevStorage) CreateLASFSMember(newMember types.LASFSMember) (string, error) {
	var entry string
	var err error
	index := s.findLASFSMember(newMember)
	if index == -1 {
		s.Members = append(s.Members, newMember)
		entry = strconv.Itoa(len(s.Members))
	} else {
		entry = strconv.Itoa(len(s.Members))
		err = errors.New("LASFSMember already created")
	}
	return entry, err
}

func (s *DevStorage) findLASFSMember(n types.LASFSMember) int {
	index := -1 // Initialize with -1 to indicate not found
	for i, m := range s.Members {
		if m.GetName() == n.GetName() {
			index = i
			break // Exit the loop once the target is found
		}
	}
	return index
}

func (s *DevStorage) GetLASFSMemberByName(n string) *types.LASFSMember {
	// TODO handle user not found index -1
	index := -1 // Initialize with -1 to indicate not found
	for i, m := range s.Members {
		if m.GetName() == n {
			index = i
			break // Exit the loop once the target is found
		}
	}
	return &s.Members[index]
}

func (s *DevStorage) GetLASFSMemberByID(id string) *types.LASFSMember {
	// TODO handle user not found index -1
	index, _ := strconv.Atoi(id)
	return &s.Members[index]
}

func (s *DevStorage) AddLASFSBallotToLASFSElection(id string, ballot types.LASFSBallot) {
	// TODO - check and handle duplicate username in LASFSBallots
	l := s.LASFSElections[id]
	l.LASFSBallots = append(l.LASFSBallots, ballot)
}
