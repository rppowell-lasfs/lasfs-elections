package types

import (
	"fmt"
	"slices"
)

type Election struct {
	Position string
	Nominees []string
}

type RawBallot struct {
	VoterName string
	Nominees  []string
}

func (r *RawBallot) MakeLASFSBallot() LASFSBallot {
	p := LASFSBallot{
		VoterName: r.VoterName,
	}
	for _, nominee := range r.Nominees {
		p.Nominees = append(p.Nominees, NomineeEntry{NomineeName: nominee, IsValid: true})
	}
	return p
}

type ElectionResult struct {
	Position       string
	ElectionStatus string
	Nominees       []string
	NomineeBuckets map[string][]LASFSBallot
}

func (e *ElectionResult) ProcessLASFSBallot(lasfsBallot LASFSBallot) {
	for !lasfsBallot.IsDead {
		nominee := lasfsBallot.GetNextNominee()
		fmt.Printf("nominee: '%s'\n", nominee)
		if len(nominee) != 0 && slices.Contains(e.Nominees, nominee) {
			fmt.Printf("adding nominee: '%s' to '%v'\n", nominee, e.NomineeBuckets[nominee])
			e.NomineeBuckets[nominee] = append(e.NomineeBuckets[nominee], lasfsBallot)
			fmt.Printf("append nominee: '%v'\n", e.NomineeBuckets[nominee])
			break
		} else {
			fmt.Printf("scratch nominee: '%s'\n", nominee)
			lasfsBallot.ScratchNominee(nominee)
		}
	}
}

func (e *ElectionResult) ProcessLASFSBallots(lasfsBallots []LASFSBallot) {
	for _, lasfsBallot := range lasfsBallots {
		e.ProcessLASFSBallot(lasfsBallot)
	}
}

type NomineeTally struct {
	Nominee string
	Count   int
}
type ElectionResultReport struct {
	Position       string
	ElectionStatus string
	Tally          []NomineeTally
}

func (e *ElectionResult) MakeElectionResultReport() ElectionResultReport {
	report := ElectionResultReport{
		Position:       e.Position,
		ElectionStatus: e.ElectionStatus,
	}
	for _, n := range e.Nominees {
		report.Tally = append(report.Tally, NomineeTally{
			Nominee: n,
			Count:   len(e.NomineeBuckets[n]),
		})
	}
	return report
}

type LASFSBallot struct {
	VoterName string
	Nominees  []NomineeEntry
	IsDead    bool
}

func (l *LASFSBallot) GetNextNominee() string {
	if l.IsDead {
		return ""
	}
	for _, nomineeEntry := range l.Nominees {
		if nomineeEntry.IsValid {
			return nomineeEntry.NomineeName
		}
	}
	return ""
}

func (l *LASFSBallot) ScratchNominee(nominee string) {
	if !l.IsDead {
		for i, nomineeEntry := range l.Nominees {
			if nomineeEntry.IsValid && nomineeEntry.NomineeName == nominee {
				l.Nominees[i].IsValid = false
			}
		}
		l.UpdateBallotStatus()
	}
}

func (l *LASFSBallot) UpdateBallotStatus() {
	if !l.IsDead {
		for _, nomineeEntry := range l.Nominees {
			if nomineeEntry.IsValid {
				return
			}
		}
		l.IsDead = true
	}
}

type NomineeEntry struct {
	NomineeName string
	IsValid     bool
}
