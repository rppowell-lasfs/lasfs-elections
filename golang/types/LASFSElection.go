package types

import "time"

type LASFSElection struct {
	ID              string
	Position        string
	Nominees        []string
	Status          string
	LASFSBallots    []LASFSBallot
	DateTimeCreated time.Time
}

func NewLASFSElection(position string, createdDateTime time.Time) *LASFSElection {
	id := makeLASFSElectionID(position, createdDateTime)
	return &LASFSElection{
		ID:              id,
		Position:        position,
		DateTimeCreated: createdDateTime,
	}
}

func makeLASFSElectionID(position string, createdDateTime time.Time) string {
	return createdDateTime.UTC().Format(time.RFC3339) + " " + position
}

func (l *LASFSElection) IsOpen() bool {
	return l.Status == "OPEN"
}

func (l *LASFSElection) IsClosed() bool {
	return l.Status == "CLOSED"
}
