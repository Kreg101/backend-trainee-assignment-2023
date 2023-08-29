package server

import (
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Storage interface for store
type Storage interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	CreateUser(int64) (int64, error)
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
	e := echo.New()

	e.POST("/segments", s.createSegment)
	e.DELETE("/segments", s.deleteSegment)
	e.POST("/users", s.createUser)
	e.PATCH("/users", s.updateUser)
	e.GET("/user", s.getUser)

	return e.Start(s.listenAddr)
}

// createSegment creates new segment
func (s *HttpServer) createSegment(c echo.Context) error {
	return nil
}

// deleteSegment deletes existing segment
func (s *HttpServer) deleteSegment(c echo.Context) error {
	return nil
}

// createUser creates new user and returns id
func (s *HttpServer) createUser(c echo.Context) error {
	var userID UserID
	err := c.Bind(&userID)
	if err != nil {
		return err
	}

	id, err := s.storage.CreateUser(userID.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, UserID{Id: id})
}

// updateUser update existing user
func (s *HttpServer) updateUser(c echo.Context) error {
	return nil
}

// getUser gets user by id
func (s *HttpServer) getUser(c echo.Context) error {
	return nil
}
