package server

import (
	"github.com/labstack/echo/v4"
	"reflect"
	"testing"
)

func TestHttpServer_Run(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpServer_addSegmentsToUser(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.addSegmentsToUser(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("addSegmentsToUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpServer_createSegment(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.createSegment(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("createSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpServer_createUser(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.createUser(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("createUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpServer_deleteSegment(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.deleteSegment(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("deleteSegment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpServer_deleteSegmentsFromUser(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.deleteSegmentsFromUser(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("deleteSegmentsFromUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpServer_getUser(t *testing.T) {
	type fields struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpServer{
				listenAddr: tt.fields.listenAddr,
				storage:    tt.fields.storage,
				logger:     tt.fields.logger,
			}
			if err := s.getUser(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("getUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		listenAddr string
		storage    Storage
		logger     *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want *HttpServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.listenAddr, tt.args.storage, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withLogging(t *testing.T) {
	type args struct {
		e      *echo.Echo
		logger *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withLogging(tt.args.e, tt.args.logger)
		})
	}
}
