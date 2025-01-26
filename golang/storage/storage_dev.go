package storage

import (
	"election/types"
	"errors"
	"fmt"
	"sort"
	"time"
)

type DevStorage struct {
	LASFSElections map[string]*types.LASFSElection
	LASFSMembers   map[string]types.LASFSMember
	LASFSMemberIDs []string
}

func NewDevStorage() *DevStorage {
	return &DevStorage{
		LASFSElections: make(map[string]*types.LASFSElection),
		LASFSMembers:   make(map[string]types.LASFSMember),
		LASFSMemberIDs: make([]string, 0),
	}
}

func (s *DevStorage) NewLASFSElection(position string, createdDateTime time.Time, nominees []string) *types.LASFSElection {
	e := types.NewLASFSElection(position, createdDateTime, nominees)
	id := e.ID
	e.DateTimeCreated = createdDateTime
	s.LASFSElections[id] = &e
	return &e
}

func (s *DevStorage) GetLASFSElectionsIDs() []string {
	keys := make([]string, 0, len(s.LASFSElections))

	for k := range s.LASFSElections {
		keys = append(keys, s.LASFSElections[k].ID)
	}
	sort.Strings(keys)
	return keys
}

func (s *DevStorage) GetLASFSElectionByID(id string) (*types.LASFSElection, error) {
	m, ok := s.LASFSElections[id]
	if !ok {
		return nil, fmt.Errorf("unable to find election '%s'", id)
	}
	return m, nil
}

func (s *DevStorage) CreateLASFSMember(username string, password string) (*types.LASFSMember, error) {
	_, err := s.GetLASFSMemberByName(username)
	if err != nil {
		newMember := types.NewLASFSMember(username, password)
		id := newMember.GetID()

		existingmember, _ := s.GetLASFSMemberByID(id)
		if existingmember != nil {
			return nil, fmt.Errorf("unable to create LASFSMember, duplicate ID '%s'", id)
		} else {
			s.LASFSMembers[id] = *newMember
			s.LASFSMemberIDs = append(s.LASFSMemberIDs, id)
			return newMember, nil
		}
	} else {
		return nil, fmt.Errorf("unable to create LASFSMember, name '%s' already exists", username)
	}
}

func (s *DevStorage) GetLASFSMemberByName(n string) (*types.LASFSMember, error) {
	for _, m := range s.LASFSMembers {
		if m.GetName() == n {
			return &m, nil
		}
	}
	return nil, errors.New("unable to find member)")
}

func (s *DevStorage) GetLASFSMemberByID(id string) (*types.LASFSMember, error) {
	m, ok := s.LASFSMembers[id]
	if ok {
		return &m, nil
	} else {
		return nil, fmt.Errorf("unable to find member by id '%s'", id)
	}
}

func (s *DevStorage) GetLASFSBallot(electionId string, memberId string) (*types.LASFSBallot, error) {
	e, err := s.GetLASFSElectionByID(electionId)
	if err != nil {
		return nil, err
	} else {
		return e.GetLASFSMemberBallot(memberId)
	}
}

func (s *DevStorage) NewLASFSBallot(memberId string, nominees []string) (*types.LASFSBallot, error) {
	member, err := s.GetLASFSMemberByID(memberId)
	if err != nil {
		return nil, err
	}
	ballot := types.NewLASFSBallot(member.ID, member.Name, nominees)
	return &ballot, nil
}

func (s *DevStorage) AddLASFSBallot(electionId string, ballot types.LASFSBallot) (*types.LASFSElection, error) {
	e, err := s.GetLASFSElectionByID(electionId)
	if err != nil {
		return nil, err // unable to find election
	}
	_, err = e.GetLASFSMemberBallot(ballot.VoterID)
	if err == nil {
		return nil, errors.New("ballot exists")
	}
	e.LASFSBallots = append(e.LASFSBallots, &ballot)
	return e, nil
}

func (s *DevStorage) AddNewLASFSBallot(election_id string, member_id string, nominees []string) (*types.LASFSElection, error) {
	member, err := s.GetLASFSMemberByID(member_id)
	if err != nil {
		return nil, err
	}
	ballot := types.NewLASFSBallot(member.ID, member.Name, nominees)
	e, err := s.GetLASFSElectionByID(election_id)
	if err != nil {
		return nil, err // unable to find election
	}

	_, err = e.GetLASFSMemberBallot(ballot.VoterID)
	if err != nil {
		return nil, errors.New("ballot exists")
	}
	e.LASFSBallots = append(e.LASFSBallots, &ballot)
	return e, nil

}

func (s *DevStorage) GetLASFSElectionResultAtCount(election_id string, count int) (*types.LASFSElectionResult, error) {
	lasfsElection, err := s.GetLASFSElectionByID(election_id)
	if err != nil {
		return nil, err
	}
	electionResult := lasfsElection.InitializeElectionResultReport()
	for i, ballot := range lasfsElection.LASFSBallots {
		if i < count {
			electionResult.ProcessLASFSBallot(*ballot)
		} else {
			break
		}
	}
	return electionResult, nil
}

func (s *DevStorage) GetLASFSElectionResult(election_id string) (*types.LASFSElectionResult, error) {
	lasfsElection, err := s.GetLASFSElectionByID(election_id)
	if err != nil {
		return nil, err
	}
	electionResult := lasfsElection.InitializeElectionResultReport()
	for _, ballot := range lasfsElection.LASFSBallots {
		electionResult.ProcessLASFSBallot(*ballot)
	}
	return electionResult, nil
}
