package storage_test

// func TestStorageNewLASFSElection(t *testing.T) {
// 	t.Run(
// 		"Storage newLASFSElection 01",
// 		func(t *testing.T) {
// 			re := regexp.MustCompile(`^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2}):(?P<MM>\d{2}):(?P<SS>\d{2})Z Test$`)

// 			s := storage.NewStorage()
// 			n := time.Now().UTC()
// 			id := s.NewLASFSElection("Test", n)

// 			if !re.MatchString(id) {
// 				t.Errorf("id is not in 'YYYY-mm-ddTHH:MM:SSZ Test' format: '%v'", id)
// 			}
// 			expectedId := n.Format(time.RFC3339) + " Test"
// 			if id != expectedId {
// 				t.Errorf("id '%v' does not match expected value '%v'", id, expectedId)
// 			}
// 		},
// 	)
// }
