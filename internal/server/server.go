package server

import (
	"errors"
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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

type Segment struct {
	Name string `json:"segment"`
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
		s.logger.Info("can't unmarshal createSegment json", zap.Error(err))
		return err
	}

	if segment.Name == "" {
		s.logger.Info("invalid segment name in createSegment", zap.Error(err))
		return errors.New("invalid segment name")
	}

	err = s.storage.CreateSegment(segment.Name)
	if err != nil {
		s.logger.Info("can't create segment", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusCreated, segment)
}

// deleteSegment deletes existing segment
func (s *HttpServer) deleteSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		s.logger.Info("can't unmarshal deleteSegment json", zap.Error(err))
		return err
	}

	if segment.Name == "" {
		s.logger.Info("invalid segment name in deleteSegment")
		return errors.New("invalid segment name")
	}

	err = s.storage.DeleteSegment(segment.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, segment)
}

// createUser creates new user and returns id
func (s *HttpServer) createUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info("can't unmarshal createUser json", zap.Error(err))
		return err
	}

	if user.Id <= 0 {
		s.logger.Info("invalid user id in createUser")
		return errors.New("invalid user id")
	}

	err = s.storage.CreateUser(user.Id)
	if err != nil {
		s.logger.Info("can't create user", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

// addSegmentsToUser add segments to existing user
func (s *HttpServer) addSegmentsToUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info("can't unmarshal addSegmentsToUser json", zap.Error(err))
		return err
	}

	if user.Id <= 0 {
		s.logger.Info("invalid user id in addSegmentsToUser")
		return errors.New("invalid user id")
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
		s.logger.Info("can't unmarshal deleteSegmentsFromUser json", zap.Error(err))
		return err
	}

	if user.Id <= 0 {
		s.logger.Info("invalid user id in deleteSegmentsFromUser")
		return errors.New("invalid user id")
	}

	err = s.storage.DeleteSegmentsFromUser(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

// getUser gets user by id
func (s *HttpServer) getUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info("can't unmarshal getUser json", zap.Error(err))
		return err
	}

	if user.Id <= 0 {
		s.logger.Info("invalid user id in getUser")
		return errors.New("invalid user id")
	}

	user, err = s.storage.GetUser(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

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
