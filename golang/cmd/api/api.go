package api

type ErrorResponsePayload struct {
	ErrorMessage string `json:"error"`
}

type CreateLASFSElectionRequestPayload struct {
	Position string   `json:"position"`
	Nominees []string `json:"nominees"`
}

type GetLASFSElectionsResponsePayload struct {
	LASFSElections []string `json:"elections"`
}

type GetLASFSElectionResponsePayload struct {
	ElectionID       string   `json:"election"`
	ElectionPosition string   `json:"position"`
	ElectionStatus   string   `json:"status"`
	Nominees         []string `json:"nominees"`
}

type CreateLASFSMemberRequestPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type CreateLASFSMemberResponsePayload struct {
	ID string `json:"id"`
}

type LASFSMemberLoginRequestPayload struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateLASFSBallotRequestPayload struct {
	Election string   `json:"election"`
	VoterID  string   `json:"id"`
	Nominees []string `json:"nominees"`
}

type GetLASFSBallotResponsePayload struct {
	Election string   `json:"election"`
	VoterID  string   `json:"id"`
	Nominees []string `json:"nominees"`
}

type PostLASFSBallotRequestPayload struct {
	Nominees []string `json:"nominees"`
}

type GetLASFSElectionResultReport struct {
	Position       string `json:"position"`
	ElectionStatus string `json:"status"`
	// TODO
	// Tally          []types.NomineeTally `json:"tally"`
}
