package types_test

import (
	"election/types"
	"reflect"
	"testing"
)

func TestNomineeTally(t *testing.T) {
	t.Run(
		"TestNomineeTally RankNomineeTally 01",
		func(t *testing.T) {
			nomineeTally := []types.NomineeTally{
				{Nominee: "Nominee01", Count: 1},
				{Nominee: "Nominee02", Count: 2},
				{Nominee: "Nominee03", Count: 3},
			}
			expectedNomineeTally := []types.NomineeTally{
				{Nominee: "Nominee03", Count: 3},
				{Nominee: "Nominee02", Count: 2},
				{Nominee: "Nominee01", Count: 1},
			}

			nomineeTally = types.RankNomineeTally(nomineeTally)
			if !reflect.DeepEqual(nomineeTally, expectedNomineeTally) {
				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", nomineeTally, expectedNomineeTally)
			}
		},
	)
}

func TestElectionResultReport(t *testing.T) {
	t.Run(
		"TestElectionResultReport RankNomineeTally 01",
		func(t *testing.T) {
			electionResult := types.LASFSElectionResult{
				Position: "Test",
				Nominees: []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: map[string][]types.LASFSBallot{
					"Nominee01": {
						{
							VoterName: "Voter001",
							Nominees: []types.NomineeEntry{
								{NomineeName: "Nominee01", IsValid: true},
								{NomineeName: "Nominee02", IsValid: true},
								{NomineeName: "Nominee03", IsValid: true},
							},
							IsDead: false,
						},
						{
							VoterName: "Voter004",
							Nominees: []types.NomineeEntry{
								{NomineeName: "Nominee01", IsValid: true},
								{NomineeName: "Nominee02", IsValid: true},
								{NomineeName: "Nominee03", IsValid: true},
							},
							IsDead: false,
						},
					},
					"Nominee02": {
						{
							VoterName: "Voter002",
							Nominees: []types.NomineeEntry{
								{NomineeName: "Nominee02", IsValid: true},
								{NomineeName: "Nominee01", IsValid: true},
								{NomineeName: "Nominee03", IsValid: true},
							},
							IsDead: false,
						},
					},
					"Nominee03": {
						{
							VoterName: "Voter003",
							Nominees: []types.NomineeEntry{
								{NomineeName: "Nominee03", IsValid: true},
								{NomineeName: "Nominee01", IsValid: true},
								{NomineeName: "Nominee02", IsValid: true},
							},
							IsDead: false,
						},
					},
				},
			}

			electionResultReport := electionResult.MakeElectionResultReport()
			expectedElectionResultReport := types.LASFSElectionResultReport{
				Position: "Test",
				Tally: []types.NomineeTally{
					{Nominee: "Nominee01", Count: 2},
					{Nominee: "Nominee02", Count: 1},
					{Nominee: "Nominee03", Count: 1},
				},
			}
			if !reflect.DeepEqual(electionResultReport, expectedElectionResultReport) {
				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", electionResultReport, expectedElectionResultReport)
			}
		},
	)
}
