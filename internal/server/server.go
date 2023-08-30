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
	CreateUser(int64) error
	AddSegmentsToUser(user db.User) error
	DeleteSegmentsFromUser(user db.User) error
	GetUser(id int64) (db.User, error)
}

type UserID struct {
	Id int64 `json:"id"`
}

type Segment struct {
	Name string `json:"segment"`
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
	e.POST("/users", s.addSegmentsToUser)
	e.DELETE("/users", s.deleteSegmentsFromUser)
	e.GET("/users", s.getUser)

	return e.Start(s.listenAddr)
}

// createSegment creates new segment
func (s *HttpServer) createSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		return err
	}

	err = s.storage.CreateSegment(segment.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, segment)
}

// deleteSegment deletes existing segment
func (s *HttpServer) deleteSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		return err
	}

	err = s.storage.DeleteSegment(segment.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, segment)
}

// createUser creates new user and returns id
func (s *HttpServer) createUser(c echo.Context) error {
	var userID UserID
	err := c.Bind(&userID)
	if err != nil {
		return err
	}

	err = s.storage.CreateUser(userID.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, userID)
}

// addSegmentsToUser add segments to existing user
func (s *HttpServer) addSegmentsToUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	err = s.storage.AddSegmentsToUser(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// deleteSegmentsFromUser deletes segments from existing user
func (s *HttpServer) deleteSegmentsFromUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	err = s.storage.DeleteSegmentsFromUser(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// getUser gets user by id
func (s *HttpServer) getUser(c echo.Context) error {
	var userID UserID
	err := c.Bind(&userID)
	if err != nil {
		return err
	}

	user, err := s.storage.GetUser(userID.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
