package types_test

import (
	"election/types"
	"reflect"
	"testing"
)

func TestLASFSBallot(t *testing.T) {

	rawBallot := types.RawBallot{
		VoterName: "Voter001",
		Nominees:  []string{"Nominee01", "Nominee02", "Nominee03"},
	}

	// r := (&types.RawBallot{
	// 	VoterName: "Voter001",
	// 	Nominees:  []string{"Nominee01", "Nominee02", "Nominee03"},
	// }).MakeLASFSBallot()
	// log.Printf("%v\n", r)

	t.Run(
		"LASFSBallot Initalize Nominee01",
		func(t *testing.T) {
			lasfsBallot := rawBallot.MakeLASFSBallot()
			expectedLASFSBallot := types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: true},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("NewLASFSBallot() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}

			nextNominee := lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee01" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee01")
			}
		},
	)
	t.Run(
		"LASFSBallot Scratch Nominee01",
		func(t *testing.T) {
			lasfsBallot := rawBallot.MakeLASFSBallot()
			var expectedLASFSBallot types.LASFSBallot
			var nextNominee string

			// Test single scratch
			lasfsBallot.ScratchNominee("Nominee01")
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee02" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee02")
			}
			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				},
				IsDead: false,
			}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("ScratchNominee() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			if lasfsBallot.IsDead {
				t.Errorf("lasfsBallot is unexpectedly dead: '%v'", lasfsBallot)
			}

			// Test double scratch
			lasfsBallot.ScratchNominee("Nominee01")
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee02" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee02")
			}
			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				},
				IsDead: false,
			}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("ScratchNominee() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			if lasfsBallot.IsDead {
				t.Errorf("lasfsBallot is unexpectedly dead: '%v'", lasfsBallot)
			}

			lasfsBallot.ScratchNominee("Nominee02")
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee03" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee03")
			}
			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: false},
					{NomineeName: "Nominee03", IsValid: true},
				},
				IsDead: false,
			}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("ScratchNominee() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			if lasfsBallot.IsDead {
				t.Errorf("lasfsBallot is unexpectedly dead: '%v'", lasfsBallot)
			}

			lasfsBallot.ScratchNominee("Nominee03")
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "")
			}
			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: false},
					{NomineeName: "Nominee03", IsValid: false},
				},
				IsDead: true,
			}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("ScratchNominee() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			if !lasfsBallot.IsDead {
				t.Errorf("lasfsBallot is unexpectedly alive '%v'", lasfsBallot)
			}

			// Scratch Nominee03 again
			lasfsBallot.ScratchNominee("Nominee03")
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "")
			}
			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: false},
					{NomineeName: "Nominee03", IsValid: false},
				},
				IsDead: true,
			}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("ScratchNominee() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			if !lasfsBallot.IsDead {
				t.Errorf("lasfsBallot is unexpectedly alive '%v'", lasfsBallot)
			}
		},
	)
}

func TestElection(t *testing.T) {
	t.Run(
		"ElectionResult ProcessLASFSBallot 01",
		func(t *testing.T) {
			electionResult := types.ElectionResult{
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

			expectedElectionResult := types.ElectionResult{
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
			expectedElectionResultReport := types.ElectionResultReport{
				Position: "Test",
				Tally: []types.NomineeTally{
					{Nominee: "Nominee001", Count: 1},
					{Nominee: "Nominee001", Count: 0},
					{Nominee: "Nominee001", Count: 0},
				},
			}
			if !reflect.DeepEqual(electionResult, expectedElectionResult) {

				t.Errorf("ProcessLASFSBallot() ElectionResultReport got '%v', expecting '%v'", electionResultReport, expectedElectionResultReport)
			}

		},
	)
	t.Run(
		"ElectionResult ProcessLASFSBallot 3 x 01",
		func(t *testing.T) {
			electionResult := types.ElectionResult{
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
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee02", IsValid: true},
						{NomineeName: "Nominee03", IsValid: true},
					},
					IsDead: false,
				},
				{
					VoterName: "Voter003",
					Nominees: []types.NomineeEntry{
						{NomineeName: "Nominee01", IsValid: true},
						{NomineeName: "Nominee02", IsValid: true},
						{NomineeName: "Nominee03", IsValid: true},
					},
					IsDead: false,
				},
			}
			electionResult.ProcessLASFSBallots(lasfsBallots)

			expectedElectionResult := types.ElectionResult{
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
							VoterName: "Voter002",
							Nominees: []types.NomineeEntry{
								{NomineeName: "Nominee01", IsValid: true},
								{NomineeName: "Nominee02", IsValid: true},
								{NomineeName: "Nominee03", IsValid: true},
							},
							IsDead: false,
						},
						{
							VoterName: "Voter003",
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
		},
	)
}
