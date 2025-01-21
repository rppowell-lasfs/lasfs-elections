package types

import "sort"

type LASFSElectionResultReport struct {
	Position       string
	ElectionStatus string
	Tally          []NomineeTally
}

func (e *LASFSElectionResult) MakeElectionResultReport() LASFSElectionResultReport {
	report := LASFSElectionResultReport{
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

func (e *LASFSElectionResultReport) IsConclusive() bool {
	sort.Slice(e.Tally, func(i, j int) bool {
		return e.Tally[i].Count > e.Tally[j].Count
	})

	if len(e.Tally) == 1 {
		return true
	} else if len(e.Tally) >= 2 && e.Tally[0].Count > e.Tally[1].Count {
		return true
	} else {
		return false
	}
}

func (e *LASFSElectionResultReport) GetWinner() string {
	sort.Slice(e.Tally, func(i, j int) bool {
		return e.Tally[i].Count > e.Tally[j].Count
	})

	if len(e.Tally) == 1 {
		return e.Tally[0].Nominee
	} else if len(e.Tally) >= 2 && e.Tally[0].Count > e.Tally[1].Count {
		return e.Tally[0].Nominee
	} else {
		return ""
	}
}
