package storage

import (
	"election/types"
	"time"
)

type StorageInterface interface {
	// NewStorage() Storage
	NewLASFSElection(position string, createdDateTime time.Time, nominees []string) *types.LASFSElection
	GetLASFSElectionsIDs() []string
	GetLASFSElectionByID(id string) (*types.LASFSElection, error)
	GetLASFSBallot(election_id string, member_id string) (*types.LASFSBallot, error)
	NewLASFSBallot(memberId string, nominees []string) (*types.LASFSBallot, error)
	AddLASFSBallot(election_id string, ballot types.LASFSBallot) (*types.LASFSElection, error)
	AddNewLASFSBallot(election_id string, member_id string, nominees []string) (*types.LASFSElection, error)
	GetLASFSElectionResult(id string) (*types.LASFSElectionResult, error)
	StorageLASFSMemberInterface
}

type StorageLASFSMemberInterface interface {
	CreateLASFSMember(username string, password string) (*types.LASFSMember, error)
	GetLASFSMemberByID(id string) (*types.LASFSMember, error)
	GetLASFSMemberByName(name string) (*types.LASFSMember, error)
}
