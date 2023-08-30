package server

import (
	"github.com/Kreg101/backend-trainee-assignment-2023/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
)

// Storage interface for store users and their segments
type Storage interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	CheckSegment(name string) (bool, error)
	CreateUser(id int64) error
	AddSegmentsToUser(user db.User) error
	DeleteSegmentsFromUser(user db.User) error
	GetUser(id int64) (db.User, error)
	CheckUser(id int64) (bool, error)
}

// Segment structure for json unmarshalling
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
	e.PATCH("/users", s.addSegmentsToUser)
	e.DELETE("/users", s.deleteSegmentsFromUser)
	e.GET("/users", s.getUser)

	return e.Start(s.listenAddr)
}

// createSegment creates new segment
func (s *HttpServer) createSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, "can't unmarshal json")
	}

	// can't create segment with empty name
	if segment.Name == "" {
		return c.JSON(http.StatusNotFound, "invalid segment name")
	}

	// check if segment already exists
	exists, err := s.storage.CheckSegment(segment.Name)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't check segment existence")
	}

	if exists {
		return c.JSON(http.StatusBadRequest, "segment already exists")
	}

	// create new segment
	err = s.storage.CreateSegment(segment.Name)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't create segment")
	}

	return c.JSON(http.StatusCreated, segment)
}

// deleteSegment deletes existing segment
func (s *HttpServer) deleteSegment(c echo.Context) error {
	var segment Segment
	err := c.Bind(&segment)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, "can't unmarshal json")
	}

	// check that segment for removal exists
	exists, err := s.storage.CheckSegment(segment.Name)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't check segment existence")
	}

	if !exists {
		return c.JSON(http.StatusNotFound, "there is no segment with that name")
	}

	// delete segment
	err = s.storage.DeleteSegment(segment.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, segment)
}

// createUser creates new user
func (s *HttpServer) createUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, "can't unmarshal json")
	}

	// can't create user with id <= 0
	if user.Id <= 0 {
		return c.JSON(http.StatusNotFound, "invalid user id")
	}

	// check if user already exists
	exists, err := s.storage.CheckUser(user.Id)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't check user existence")
	}

	if exists {
		return c.JSON(http.StatusBadRequest, "user already exists")
	}

	// create new user
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
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusBadRequest, "can't unmarshal json")
	}

	// check that user exists
	exists, err := s.storage.CheckUser(user.Id)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't check user existence")
	}

	if !exists {
		return c.JSON(http.StatusNotFound, "user doesn't exist")
	}

	// correctUser is a user with segments that already exists
	// we should add only existing segments
	correctUser := db.User{
		Id:         user.Id,
		Segments:   make([]string, 0),
		ActiveTime: user.ActiveTime,
	}
	for _, name := range user.Segments {
		exists, err = s.storage.CheckSegment(name)
		if err != nil {
			s.logger.Info(zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "can't check segment existence")
		}
		if exists {
			correctUser.Segments = append(correctUser.Segments, name)
		}
	}

	// add existing segments to user
	err = s.storage.AddSegmentsToUser(correctUser)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, correctUser)
}

// deleteSegmentsFromUser deletes segments from existing user
func (s *HttpServer) deleteSegmentsFromUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info("can't unmarshal deleteSegmentsFromUser json", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "can't unmarshal json")
	}

	// check that user exists
	exists, err := s.storage.CheckUser(user.Id)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't check user existence")
	}

	if !exists {
		return c.JSON(http.StatusNotFound, "user doesn't exist")
	}

	// correctUser is a user with segments that already exists
	// we should delete only existing segments
	correctUser := db.User{
		Id:         user.Id,
		Segments:   make([]string, 0),
		ActiveTime: user.ActiveTime,
	}
	for _, name := range user.Segments {
		exists, err = s.storage.CheckSegment(name)
		if err != nil {
			s.logger.Info(zap.Error(err))
			return c.JSON(http.StatusInternalServerError, "can't check segment existence")
		}
		if exists {
			correctUser.Segments = append(correctUser.Segments, name)
		}
	}

	// delete existing segments
	err = s.storage.DeleteSegmentsFromUser(correctUser)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, correctUser)
}

// getUser gets user from database by it's id
func (s *HttpServer) getUser(c echo.Context) error {
	var user db.User
	err := c.Bind(&user)
	if err != nil {
		s.logger.Info("can't unmarshal getUser json", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "can't unmarshal json")
	}

	// check that user exists
	exists, err := s.storage.CheckUser(user.Id)
	if err != nil {
		s.logger.Info(zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "can't check user existence")
	}

	if !exists {
		return c.JSON(http.StatusNotFound, "user doesn't exist")
	}

	// get user's segments
	user, err = s.storage.GetUser(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
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
