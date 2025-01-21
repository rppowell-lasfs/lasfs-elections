package storage

import (
	"election/types"
	"time"
)

// type Storage interface {
// 	// NewStorage() Storage
// 	NewLASFSElection(position string, createdDateTime time.Time) string
// 	GetLASFSElectionsIDs() []string
// 	GetLASFSElectionByID(id string) types.LASFSElection
// 	CreateLASFSMember(newMember types.LASFSMember) (string, error)
// 	GetLASFSMemberByID(id string) types.LASFSMember
// 	GetLASFSMemberByName(name string) types.LASFSMember
// 	AddLASFSBallotToLASFSElection(id string, ballot types.LASFSBallot)
// }

type StorageInterface interface {
	// NewStorage() Storage
	NewLASFSElection(position string, createdDateTime time.Time) string
	GetLASFSElectionsIDs() []string
	GetLASFSElectionByID(id string) types.LASFSElection
	StorageLASFSMemberInterface
	AddLASFSBallotToLASFSElection(id string, ballot types.LASFSBallot)
}

type StorageLASFSMemberInterface interface {
	CreateLASFSMember(newMember types.LASFSMember) (string, error)
	GetLASFSMemberByID(id string) *types.LASFSMember
	GetLASFSMemberByName(name string) *types.LASFSMember
}
