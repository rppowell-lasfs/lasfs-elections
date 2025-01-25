package types_test

import (
	"election/types"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLASFSBallotFromRawBallot(t *testing.T) {
	rawBallot := types.RawBallot{
		VoterName: "Voter001",
		Nominees:  []string{"Nominee01", "Nominee02", "Nominee03"},
	}

	lasfsBallot := rawBallot.MakeLASFSBallot()
	expectedLASFSBallot := types.LASFSBallot{
		VoterID:   "Voter001",
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
}

func TestLASFSBallotFunctions(t *testing.T) {

	rawBallot := types.RawBallot{
		VoterName: "Voter001",
		Nominees:  []string{"Nominee01", "Nominee02", "Nominee03"},
	}

	t.Run(
		"LASFSBallot from RawBallot 1",
		func(t *testing.T) {
			lasfsBallot := rawBallot.MakeLASFSBallot()
			expectedLASFSBallot := types.LASFSBallot{
				VoterID:   "Voter001",
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
				VoterID:   "Voter001",
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
				VoterID:   "Voter001",
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
				VoterID:   "Voter001",
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
				VoterID:   "Voter001",
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
				VoterID:   "Voter001",
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
func TestLASFSBallotScratchFunctions(t *testing.T) {
	t.Run(
		"LASFSBallot Scratch One",
		func(t *testing.T) {
			lasfsBallot := types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: true},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			lasfsBallot.ScratchNominee("Nominee01")
			expectedLASFSBallot := types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("NewLASFSBallot() got '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			nextNominee := lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee02" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee02")
			}
		},
	)
	t.Run(
		"LASFSBallot Scratch All Three",
		func(t *testing.T) {
			lasfsBallot := types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: true},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			var expectedLASFSBallot types.LASFSBallot
			var nextNominee string

			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: true},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("LASFSBallot is '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee01" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee01")
			}

			lasfsBallot.ScratchNominee("Nominee01")

			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: true},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("LASFSBallot is '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee02" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee02")
			}
			if lasfsBallot.IsDead {
				t.Errorf("lasfsBallot unexpectedly dead")
			}

			lasfsBallot.ScratchNominee("Nominee02")
			expectedLASFSBallot = types.LASFSBallot{
				VoterName: "Voter001",
				Nominees: []types.NomineeEntry{
					{NomineeName: "Nominee01", IsValid: false},
					{NomineeName: "Nominee02", IsValid: false},
					{NomineeName: "Nominee03", IsValid: true},
				}}
			if !reflect.DeepEqual(lasfsBallot, expectedLASFSBallot) {
				t.Errorf("LASFSBallot is '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "Nominee03" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "Nominee03")
			}
			if lasfsBallot.IsDead {
				t.Errorf("lasfsBallot unexpectedly dead")
			}

			lasfsBallot.ScratchNominee("Nominee03")
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
				t.Errorf("LASFSBallot is '%v', expecting '%v'", lasfsBallot, expectedLASFSBallot)
			}
			nextNominee = lasfsBallot.GetNextNominee()
			if nextNominee != "" {
				t.Errorf("GetNextNominee() got '%v', expecting '%v'", nextNominee, "")
			}
			if !lasfsBallot.IsDead {
				t.Errorf("lasfsBallot unexpectedly alive")
			}

		},
	)
}

func TestLASFSBallotNomineeVotes(t *testing.T) {
	lasfsBallot := types.LASFSBallot{
		VoterName: "Voter001",
		Nominees: []types.NomineeEntry{
			{NomineeName: "Nominee01", IsValid: true},
			{NomineeName: "Nominee02", IsValid: true},
			{NomineeName: "Nominee03", IsValid: true},
		}}

	nomineeVotes := lasfsBallot.NomineeVotes()
	assert.Equal(t, []string{"Nominee01", "Nominee02", "Nominee03"}, nomineeVotes)

}
