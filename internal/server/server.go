package server

import (
	"github.com/gocarina/gocsv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// Storage interface for store users and their segments
type Storage interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	CreateUser(id int64) error
	AddSegmentsToUser(user User) error
	DeleteSegmentsFromUser(user User) error
	GetUser(id int64) (*User, error)
	GetUserHistory(user User) ([]TimeUser, error)
}

// HttpServer connects database with http requests
type HttpServer struct {
	listenAddr string
	storage    Storage
	logger     *zap.SugaredLogger
}

// NewServer creates new HttpServer
func NewServer(listenAddr string, storage Storage, logger *zap.SugaredLogger) *HttpServer {
	return &HttpServer{
		listenAddr: listenAddr,
		storage:    storage,
		logger:     logger,
	}
}

// Run configures the server and starts it
func (s *HttpServer) Run() error {
	e := echo.New()

	withLogging(e, s.logger)

	e.POST("/segments", s.createSegment)
	e.DELETE("/segments", s.deleteSegment)
	e.POST("/users", s.createUser)
	e.PATCH("/users", s.addSegmentsToUser)
	e.DELETE("/users", s.deleteSegmentsFromUser)
	e.GET("/users/:id", s.getUser)
	e.GET("/users/:id/history", s.userHistory)

	return e.Start(s.listenAddr)
}

// createSegment creates new segment
func (s *HttpServer) createSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"can't unmarshal json"})
	}

	if segment.Name == "" {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"invalid request"})
	}

	// create new segment
	err = s.storage.CreateSegment(segment.Name)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't create segment"})
	}

	return c.JSON(http.StatusCreated, segment)
}

// deleteSegment deletes existing segment
func (s *HttpServer) deleteSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"can't unmarshal json"})
	}

	// delete segment
	err = s.storage.DeleteSegment(segment.Name)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't delete segment"})
	}

	return c.JSON(http.StatusOK, segment)
}

// createUser creates new user
func (s *HttpServer) createUser(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"can't unmarshal json"})
	}

	// create new user
	err = s.storage.CreateUser(user.Id)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

// addSegmentsToUser add segments to existing user
func (s *HttpServer) addSegmentsToUser(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"can't unmarshal json"})
	}

	// add segments to user
	err = s.storage.AddSegmentsToUser(user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't add segments to user"})
	}

	return c.JSON(http.StatusOK, user)
}

// deleteSegmentsFromUser deletes segments from existing user
func (s *HttpServer) deleteSegmentsFromUser(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"can't unmarshal json"})
	}

	// delete segments from user
	err = s.storage.DeleteSegmentsFromUser(user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't delete segments from user"})
	}

	return c.JSON(http.StatusOK, user)
}

// getUser gets user from database by it's id
func (s *HttpServer) getUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"invalid user id"})
	}

	// it means db error
	user, err := s.storage.GetUser(int64(id))
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't get user"})
	}

	// it means user doesn't exist
	if user == nil {
		return c.JSON(http.StatusNotFound, Msg{"user doesn't exists"})
	}

	return c.JSON(http.StatusOK, *user)
}

// userHistory returns csv file with history of user-segment relationships
// in specified year and month
func (s *HttpServer) userHistory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"invalid user id"})
	}

	year, err := strconv.Atoi(c.QueryParam("year"))
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"invalid year"})
	}

	month, err := strconv.Atoi(c.QueryParam("month"))
	if err != nil || month < 1 || month > 12 {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, Msg{"invalid month"})
	}

	user := User{Id: int64(id), Year: year, Month: month}
	history, err := s.storage.GetUserHistory(user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't get user's history"})
	}

	c.Response().Unwrap().Header().Set("Content-Disposition", "attachment; filename=test.csv")
	err = gocsv.Marshal(history, c.Response().Unwrap())
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Msg{"can't marshal history to json"})
	}

	return nil
}

// withLogging is middleware for logging
func withLogging(e *echo.Echo, logger *zap.SugaredLogger) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("Method", v.Method),
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
}
