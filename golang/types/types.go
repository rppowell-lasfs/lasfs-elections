package types

import (
	"sort"
)

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

type NomineeTally struct {
	Nominee string
	Count   int
}

func RankNomineeTally(nomineeTally []NomineeTally) []NomineeTally {
	sort.Slice(nomineeTally, func(i, j int) bool {
		return nomineeTally[i].Count > nomineeTally[j].Count
	})
	return nomineeTally
}

type NomineeEntry struct {
	NomineeName string
	IsValid     bool
}
