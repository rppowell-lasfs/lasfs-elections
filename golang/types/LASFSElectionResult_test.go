package types_test

import (
	"election/types"
	"reflect"
	"testing"
)

func TestElectionResultScenarios01HappyPath(t *testing.T) {
	t.Run(
		"ElectionResult ProcessLASFSBallot 01 Single Vote And Result",
		func(t *testing.T) {
			electionResult := types.LASFSElectionResult{
				Position:       "Test",
				Nominees:       []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: make(map[string][]types.LASFSBallot),
			}
			lasfsBallots := []types.LASFSBallot{
				{
					VoterName: "Voter001",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee02", IsValid: true},
						{NomineeName: "Nominee03", IsValid: true},
					},
					IsDead: false,
				},
			}
			electionResult.ProcessLASFSBallot(lasfsBallots[0])

			expectedElectionResult := types.LASFSElectionResult{
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
					},
				},
			}

			if !reflect.DeepEqual(electionResult, expectedElectionResult) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionResult, expectedElectionResult)
			}

			electionResultReport := electionResult.MakeElectionResultReport()
			expectedElectionResultReport := types.LASFSElectionResultReport{
				Position: "Test",
				Tally: []types.NomineeTally{
					{Nominee: "Nominee01", Count: 1},
					{Nominee: "Nominee02", Count: 0},
					{Nominee: "Nominee03", Count: 0},
				},
			}
			if !reflect.DeepEqual(electionResultReport, expectedElectionResultReport) {
				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", electionResultReport, expectedElectionResultReport)
			}
			electionIsConclusive := electionResultReport.IsConclusive()
			if electionIsConclusive != true {
				t.Errorf("ElectionResultReport found IsConclusive '%v', expecting '%v'", electionIsConclusive, true)
			}
			winner := electionResultReport.GetWinner()
			if winner != "Nominee01" {
				t.Errorf("ElectionResultReport winner '%v', expecting '%v'", winner, "Nominee01")
			}
		},
	)
	t.Run(
		"ElectionResult ProcessLASFSBallot 3 x 01 inconclusive",
		func(t *testing.T) {
			electionResult := types.LASFSElectionResult{
				Position:       "Test",
				Nominees:       []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: make(map[string][]types.LASFSBallot),
			}
			lasfsBallots := []types.LASFSBallot{
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
					VoterName: "Voter002",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee02", IsValid: true},
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee03", IsValid: true},
					},
					IsDead: false,
				},
				{
					VoterName: "Voter003",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee03", IsValid: true},
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee02", IsValid: true},
					},
					IsDead: false,
				},
			}
			electionResult.ProcessLASFSBallots(lasfsBallots)

			expectedElectionResult := types.LASFSElectionResult{
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

			if !reflect.DeepEqual(electionResult, expectedElectionResult) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionResult, expectedElectionResult)
			}

			electionResultReport := electionResult.MakeElectionResultReport()
			expectedElectionResultReport := types.LASFSElectionResultReport{
				Position: "Test",
				Tally: []types.NomineeTally{
					{Nominee: "Nominee01", Count: 1},
					{Nominee: "Nominee02", Count: 1},
					{Nominee: "Nominee03", Count: 1},
				},
			}
			if !reflect.DeepEqual(electionResultReport, expectedElectionResultReport) {

				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", electionResultReport, expectedElectionResultReport)
			}
			electionIsConclusive := electionResultReport.IsConclusive()
			if electionIsConclusive != false {
				t.Errorf("ElectionResultReport found IsConclusive '%v', expecting '%v'", electionIsConclusive, false)
			}
			winner := electionResultReport.GetWinner()
			if winner != "" {
				t.Errorf("ElectionResultReport winner '%v', expecting '%v'", winner, "")
			}

		},
	)
}

func TestElectionResultScenarios02(t *testing.T) {
	t.Run(
		"TestElectionScenarios02 ElectionResult 4 votes",
		func(t *testing.T) {
			electionResult := types.LASFSElectionResult{
				Position:       "Test",
				Nominees:       []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: make(map[string][]types.LASFSBallot),
			}
			lasfsBallots := []types.LASFSBallot{
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
					VoterName: "Voter002",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee02", IsValid: true},
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee03", IsValid: true},
					},
					IsDead: false,
				},
				{
					VoterName: "Voter003",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee03", IsValid: true},
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee02", IsValid: true},
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
			}
			electionResult.ProcessLASFSBallots(lasfsBallots)

			expectedElectionResult := types.LASFSElectionResult{
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

			if !reflect.DeepEqual(electionResult, expectedElectionResult) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionResult, expectedElectionResult)
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
			electionIsConclusive := electionResultReport.IsConclusive()
			if electionIsConclusive != true {
				t.Errorf("ElectionResultReport found IsConclusive '%v', expecting '%v'", electionIsConclusive, true)
			}
			winner := electionResultReport.GetWinner()
			if winner != "Nominee01" {
				t.Errorf("ElectionResultReport winner '%v', expecting '%v'", winner, "Nominee01")
			}

		},
	)
}

func TestElectionResultScenarios03InvalidNominees(t *testing.T) {
	t.Run(
		"ElectionResult ProcessLASFSBallot With 1 Invalid Nominee",
		func(t *testing.T) {
			electionResult := types.LASFSElectionResult{
				Position:       "Test",
				Nominees:       []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: make(map[string][]types.LASFSBallot),
			}
			lasfsBallots := []types.LASFSBallot{
				{
					VoterName: "Voter001",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee04", IsValid: true},
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee02", IsValid: true},
					},
					IsDead: false,
				},
			}
			electionResult.ProcessLASFSBallot(lasfsBallots[0])

			expectedElectionResult := types.LASFSElectionResult{
				Position: "Test",
				Nominees: []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: map[string][]types.LASFSBallot{
					"Nominee01": {
						{
							VoterName: "Voter001",
							Nominees: []types.NomineeEntry{
								{NomineeName: "Nominee04", IsValid: false},
								{NomineeName: "Nominee01", IsValid: true},
								{NomineeName: "Nominee02", IsValid: true},
							},
							IsDead: false,
						},
					},
				},
			}

			if !reflect.DeepEqual(electionResult, expectedElectionResult) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionResult, expectedElectionResult)
			}

			electionResultReport := electionResult.MakeElectionResultReport()
			expectedElectionResultReport := types.LASFSElectionResultReport{
				Position: "Test",
				Tally: []types.NomineeTally{
					{Nominee: "Nominee01", Count: 1},
					{Nominee: "Nominee02", Count: 0},
					{Nominee: "Nominee03", Count: 0},
				},
			}
			if !reflect.DeepEqual(electionResultReport, expectedElectionResultReport) {
				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", electionResultReport, expectedElectionResultReport)
			}
			electionIsConclusive := electionResultReport.IsConclusive()
			if electionIsConclusive != true {
				t.Errorf("ElectionResultReport found IsConclusive '%v', expecting '%v'", electionIsConclusive, true)
			}
			winner := electionResultReport.GetWinner()
			if winner != "Nominee01" {
				t.Errorf("ElectionResultReport winner '%v', expecting '%v'", winner, "Nominee01")
			}

		},
	)
	t.Run(
		"ElectionResult ProcessLASFSBallot With All Invalid Nominees",
		func(t *testing.T) {
			electionResult := types.LASFSElectionResult{
				Position:       "Test",
				Nominees:       []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: make(map[string][]types.LASFSBallot),
			}
			lasfsBallots := []types.LASFSBallot{
				{
					VoterName: "Voter001",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee04", IsValid: true},
						{NomineeName: "Nominee05", IsValid: true},
						{NomineeName: "Nominee06", IsValid: true},
					},
					IsDead: false,
				},
			}

			electionResult.ProcessLASFSBallot(lasfsBallots[0])

			expectedElectionResult := types.LASFSElectionResult{
				Position:       "Test",
				Nominees:       []string{"Nominee01", "Nominee02", "Nominee03"},
				NomineeBuckets: make(map[string][]types.LASFSBallot),
				DeadBallots: []types.LASFSBallot{
					{
						VoterName: "Voter001",
						Nominees: []types.NomineeEntry{
							{NomineeName: "Nominee04", IsValid: false},
							{NomineeName: "Nominee05", IsValid: false},
							{NomineeName: "Nominee06", IsValid: false},
						},
						IsDead: true,
					},
				},
			}

			if !reflect.DeepEqual(electionResult, expectedElectionResult) {
				t.Errorf("ProcessLASFSBallot() got '%v', expecting '%v'", electionResult, expectedElectionResult)
			}

			electionResultReport := electionResult.MakeElectionResultReport()
			expectedElectionResultReport := types.LASFSElectionResultReport{
				Position: "Test",
				Tally: []types.NomineeTally{
					{Nominee: "Nominee01", Count: 0},
					{Nominee: "Nominee02", Count: 0},
					{Nominee: "Nominee03", Count: 0},
				},
			}
			if !reflect.DeepEqual(electionResultReport, expectedElectionResultReport) {
				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", electionResultReport, expectedElectionResultReport)
			}
			electionIsConclusive := electionResultReport.IsConclusive()
			if electionIsConclusive != false {
				t.Errorf("ElectionResultReport found IsConclusive '%v', expecting '%v'", electionIsConclusive, false)
			}
			winner := electionResultReport.GetWinner()
			if winner != "" {
				t.Errorf("ElectionResultReport winner '%v', expecting '%v'", winner, "")
			}

		},
	)
}
