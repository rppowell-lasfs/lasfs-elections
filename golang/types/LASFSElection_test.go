package types_test

import (
	"election/types"
	"regexp"
	"testing"
	"time"
)

func TestLASFSElectionIDs(t *testing.T) {

	reString := `^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2})-(?P<MM>\d{2})-(?P<SS>\d{2})\.(?P<nnnnnnnnn>\d{9})Z`

	re1 := regexp.MustCompile(reString + `_Test1$`)
	n1 := time.Now().UTC()
	e1 := types.NewLASFSElection("Test1", n1, nil)

	if !re1.MatchString(e1.ID) {
		t.Errorf("id is not in 'YYYY-mm-ddTHH-MM-SS.nnnnnnnnnZ_Test1' format: '%v'", e1.ID)
	}
	expectedId1 := n1.Format(`2006-01-02T15-04-05.000000000Z0700`) + "_Test1"
	if e1.ID != expectedId1 {
		t.Errorf("id '%v' does not match expected value '%v'", e1.ID, expectedId1)
	}

	re2 := regexp.MustCompile(reString + `_Test_2$`)
	n2 := time.Now().UTC()
	e2 := types.NewLASFSElection("Test 2", n2, nil)

	if !re2.MatchString(e2.ID) {
		t.Errorf("id is not in 'YYYY-mm-ddTHH-MM-SS.nnnnnnnnnZ_Test1' format: '%v'", e2.ID)
	}
	expectedId2 := n2.Format(`2006-01-02T15-04-05.000000000Z0700`) + "_Test_2"
	if e2.ID != expectedId2 {
		t.Errorf("id '%v' does not match expected value '%v'", e2.ID, expectedId2)
	}

	re3 := regexp.MustCompile(reString + `_Test_3$`)
	n3 := time.Now().UTC()
	e3 := types.NewLASFSElection("Test_3", n3, nil)

	if !re3.MatchString(e3.ID) {
		t.Errorf("id is not in 'YYYY-mm-ddTHH-MM-SSZ_Test3' format: '%v'", e3.ID)
	}
	expectedId3 := n3.Format(`2006-01-02T15-04-05.000000000Z0700`) + "_Test_3"
	if e3.ID != expectedId3 {
		t.Errorf("id '%v' does not match expected value '%v'", e3.ID, expectedId3)
	}
}
