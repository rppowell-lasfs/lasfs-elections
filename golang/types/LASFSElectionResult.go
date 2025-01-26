package types

import (
	"slices"
)

type LASFSElectionResult struct {
	ElectionID     string
	Position       string
	ElectionStatus string
	BallotCount    int
	Nominees       []string
	NomineeBuckets map[string][]LASFSBallot
	DeadBallots    []LASFSBallot
}

func (e *LASFSElectionResult) ProcessLASFSBallot(lasfsBallot LASFSBallot) {
	for !lasfsBallot.IsDead {
		nominee := lasfsBallot.GetNextNominee()
		if len(nominee) != 0 && slices.Contains(e.Nominees, nominee) {
			// nominee found, add to bucket
			// fmt.Printf("Nomine found, adding to '%s' bucket '%v'", nominee, lasfsBallot)
			e.NomineeBuckets[nominee] = append(e.NomineeBuckets[nominee], lasfsBallot)
			break
		} else {
			// nominee not found, scratch off nominee from ballot, check ballot again
			lasfsBallot.ScratchNominee(nominee)
		}
	}
	if lasfsBallot.IsDead {
		// fmt.Printf("Append DeadBallots '%v'\n", lasfsBallot)
		e.DeadBallots = append(e.DeadBallots, lasfsBallot)
	}
}

func (e *LASFSElectionResult) ProcessLASFSBallots(lasfsBallots []LASFSBallot) {
	for _, lasfsBallot := range lasfsBallots {
		e.ProcessLASFSBallot(lasfsBallot)
	}
}
