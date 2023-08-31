package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// Storage interface for store users and their segments
type Storage interface {
	CreateSegment(segment Segment) error
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
