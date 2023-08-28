package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type HttpServer struct {
	listenAddr string
}

func NewServer(listenAddr string) *HttpServer {
	return &HttpServer{
		listenAddr: listenAddr,
	}
}

// Start configures the server and starts it
func (s *HttpServer) Start() {
	router := chi.NewRouter()

	router.Post("/segment", createSegment)
	router.Delete("/segment", deleteSegment)
	router.Post("/user", updateUser)
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

// updateUser creates new user or update existing user
func updateUser(w http.ResponseWriter, r *http.Request) {

}

// getUser gets user by id
func getUser(w http.ResponseWriter, r *http.Request) {

}
