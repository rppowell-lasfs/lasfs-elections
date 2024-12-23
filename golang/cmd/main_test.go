package main

import (
	"testing"
)

func TestOne(t *testing.T) {
	t.Run(
		"SubmittedBallot Test",
		func(t *testing.T) {
		},
	)
	// t.Run(
	// 	"test one",
	// 	func(t *testing.T) {
	// 		testElection := types.Election{
	// 			Position: "test",
	// 			Nominees: []string{"George", "Karl", "Bob"},
	// 		}
	// 		fmt.Printf("testElection: %s", testElection.Position)
	// 		fmt.Printf("testElection: %s", testElection.Nominees)

	// 		submittedBallot001 := types.SubmittedBallot{
	// 			VoterName: "one",
	// 			Nominees:  []string{"George"},
	// 		}

	// 		submittedBallots := []types.SubmittedBallot{submittedBallot001}

	// 		electionResults, err := utils.RunElection(testElection, submittedBallots)
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}
	// 		if electionResults == nil {
	// 			t.Errorf("expected ElectionResults, got 'nil'")
	// 		}
	// 	},
	// )

}
