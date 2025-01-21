package api

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// type APIServer struct {
// 	addr string
// 	db   *sql.DB
// }

// func NewAPIServer(addr string, db *sql.DB) *APIServer {
// 	return &APIServer{
// 		addr: addr,
// 		db:   db,
// 	}
// }

// func (s *APIServer) Run() error {
// 	router := mux.NewRouter()
// 	subrouter := router.PathPrefix("/api/v1").Subrouter()

// 	userStore := user.NewStore(s.db)
// 	userHandler := user.NewHandler(userStore)
// 	userHandler.RegisterRoutes(subrouter)

// 	log.Println("Listening on", s.addr)

// 	return http.ListenAndServe(s.addr, router)
// }

type CreateLASFSElectionPayload struct {
	Position string   `json:"position"`
	Nominees []string `json:"nominees"`
}

type UpdateLASFSElectionPayload struct {
	ID    int    `json:"id"`
	State string `json:"state"`
}

type CreateLASFSBallot struct {
	Election  string   `json:"election"`
	VoterName int      `json:"votername"`
	Nominees  []string `json:"nominees"`
}

type RetrieveLASFSElectionResultReport struct {
	Position       string `json:"position"`
	ElectionStatus string `json:"status"`
	// TODO
	// Tally          []types.NomineeTally `json:"tally"`
}
