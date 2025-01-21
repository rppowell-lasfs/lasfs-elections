package types

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
