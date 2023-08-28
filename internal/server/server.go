package server

import (
	"encoding/json"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Storage interface for store
type Storage interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	CreateUser() (int64, error)
	UpdateUser(user db.User) error
	GetUser(id int64) (*db.User, error)
}

type UserID struct {
	Id int64 `json:"id"`
}

// HttpServer connects database with http requests
type HttpServer struct {
	listenAddr string
	storage    Storage
}

// NewServer creates new HttpServer
func NewServer(listenAddr string, storage Storage) *HttpServer {
	return &HttpServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

// Run configures the server and starts it
func (s *HttpServer) Run() error {
	router := chi.NewRouter()

	router.Post("/segments", s.createSegment)
	router.Delete("/segments", s.deleteSegment)
	router.Post("/users", s.createUser)
	router.Patch("/users", s.updateUser)
	router.Get("/user", s.getUser)

	return http.ListenAndServe(s.listenAddr, router)
}

// createSegment creates new segment
func (s *HttpServer) createSegment(w http.ResponseWriter, r *http.Request) {

}

// deleteSegment deletes existing segment
func (s *HttpServer) deleteSegment(w http.ResponseWriter, r *http.Request) {

}

// createUser creates new user
func (s *HttpServer) createUser(w http.ResponseWriter, r *http.Request) {
	id, err := s.storage.CreateUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't create new user"))
		return
	}

	userID := UserID{Id: id}
	err = json.NewEncoder(w).Encode(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't return new user"))
		return
	}
}

// updateUser update existing user
func (s *HttpServer) updateUser(w http.ResponseWriter, r *http.Request) {

}

// getUser gets user by id
func (s *HttpServer) getUser(w http.ResponseWriter, r *http.Request) {

}
