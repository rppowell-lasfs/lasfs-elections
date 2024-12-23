package main

import (
	"election/types"
	"fmt"
)

func main() {
	fmt.Println("Hello World")

	rawBallot := types.RawBallot{
		VoterName: "Voter001",
		Nominees:  []string{"Nominee01"},
	}
	fmt.Println("submittedBallot: %z", rawBallot)

	lasfsBallot := rawBallot.MakeLASFSBallot()

	fmt.Println("processedBallot: %z", lasfsBallot)

	// testElection := types.Election{
	// 	Position: "test",
	// 	Nominees: []string{"George", "Karl", "Bob"},
	// }
	// fmt.Printf("testElection: %s", testElection.Position)
	// fmt.Printf("testElection: %s", testElection.Nominees)

	// submittedBallot001 := types.SubmittedBallot{
	// 	VoterName: "one",
	// 	Nominees:  []string{"George"},
	// }

	// submittedBallots := []types.SubmittedBallot{submittedBallot001}

	// electionResults, err := utils.RunElection(testElection, submittedBallots)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("electionResults: %z", electionResults)
}
