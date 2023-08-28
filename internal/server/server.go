package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// HttpServer connects database with http requests
type HttpServer struct {
	listenAddr string
}

// NewServer creates new HttpServer
func NewServer(listenAddr string) *HttpServer {
	return &HttpServer{
		listenAddr: listenAddr,
	}
}

// Start configures the server and starts it
func (s *HttpServer) Start() {
	router := chi.NewRouter()

	router.Post("/segments", createSegment)
	router.Delete("/segments", deleteSegment)
	router.Post("/users", createUser)
	router.Patch("/users", updateUser)
	router.Get("/user", getUser)

	// TODO: handle errors
	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		return
	}
}

// createSegment creates new segment
func createSegment(w http.ResponseWriter, r *http.Request) {

}

// deleteSegment deletes existing segment
func deleteSegment(w http.ResponseWriter, r *http.Request) {

}

// createUser creates new user
func createUser(w http.ResponseWriter, r *http.Request) {

}

// updateUser update existing user
func updateUser(w http.ResponseWriter, r *http.Request) {

}

// getUser gets user by id
func getUser(w http.ResponseWriter, r *http.Request) {

}
