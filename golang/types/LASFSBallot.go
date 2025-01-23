package types

type LASFSBallot struct {
	VoterID   string
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

func (l *LASFSBallot) NomineeVotes() []string {
	nominees := make([]string, 0)
	for _, n := range l.Nominees {
		nominees = append(nominees, n.NomineeName)
	}
	return nominees
}

func NewLASFSBallot(voterId string, VoterName string, Nominees []string) LASFSBallot {
	p := LASFSBallot{
		VoterID:   voterId,
		VoterName: VoterName,
	}
	for _, nominee := range Nominees {
		p.Nominees = append(p.Nominees, NomineeEntry{NomineeName: nominee, IsValid: true})
	}
	return p
}
