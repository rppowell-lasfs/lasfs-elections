package types

import (
	"election/utils"
	"fmt"
	"time"
)

type LASFSElection struct {
	ID              string
	Position        string
	Nominees        []string
	Status          string
	LASFSBallots    []*LASFSBallot
	DateTimeCreated time.Time
}

func NewLASFSElection(position string, createdDateTime time.Time, nominees []string) LASFSElection {
	id := makeLASFSElectionID(position, createdDateTime)
	return LASFSElection{
		ID:              id,
		Position:        position,
		DateTimeCreated: createdDateTime,
		Nominees:        nominees,
		LASFSBallots:    make([]*LASFSBallot, 0),
	}
}

func makeLASFSElectionID(position string, createdDateTime time.Time) string {
	// time.RFC3339 `2006-01-02T15:04:05Z07:00`
	// `2006-01-02T15-04-05.000000000Z0700`
	position = utils.LASFSElectionIDFromString(position)
	return createdDateTime.UTC().Format(`2006-01-02T15-04-05.000000000Z0700`) + "_" + position
}

func (l *LASFSElection) IsOpen() bool {
	return l.Status == "OPEN"
}

func (l *LASFSElection) IsClosed() bool {
	return l.Status == "CLOSED"
}

func (e *LASFSElection) GetLASFSMemberBallot(memberId string) (b *LASFSBallot, err error) {
	for _, b := range e.LASFSBallots {
		if b.VoterID == memberId {
			return b, nil
		}
	}
	return nil, fmt.Errorf("unable to find memberId '%s'", memberId)
}
