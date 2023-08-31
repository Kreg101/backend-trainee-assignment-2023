package server

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpServer_createSegment(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		request           string
		segmentName       string
		wantToCallStorage bool
		err               error
		response          response
	}{
		{
			name:              "OK",
			request:           "{\"segment\":\"a\"}",
			segmentName:       "a",
			wantToCallStorage: true,
			err:               nil,
			response: response{
				code: http.StatusCreated,
				body: "{\"segment\":\"a\"}\n",
			},
		},
		{
			name:              "Invalid json #1",
			request:           "{\"name\":\"a\"}",
			segmentName:       "",
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"invalid request\"}\n",
			},
		},
		{
			name:              "Invalid json #2",
			request:           "{\"name\":\"a",
			segmentName:       "",
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"can't unmarshal json\"}\n",
			},
		},
		{
			name:              "internal error",
			request:           "{\"segment\":\"a\"}",
			segmentName:       "a",
			wantToCallStorage: true,
			err:               errors.New(""),
			response: response{
				code: http.StatusInternalServerError,
				body: "{\"text\":\"can't create segment\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().CreateSegment(tc.segmentName).Return(tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/segments", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.createSegment(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}

func TestHttpServer_deleteSegment(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		request           string
		segmentName       string
		wantToCallStorage bool
		err               error
		response          response
	}{
		{
			name:              "OK",
			request:           "{\"segment\":\"a\"}",
			segmentName:       "a",
			wantToCallStorage: true,
			err:               nil,
			response: response{
				code: http.StatusOK,
				body: "{\"segment\":\"a\"}\n",
			},
		},
		{
			name:              "Invalid json",
			request:           "{",
			segmentName:       "",
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"can't unmarshal json\"}\n",
			},
		},
		{
			name:              "internal error",
			request:           "{\"segment\":\"a\"}",
			segmentName:       "a",
			wantToCallStorage: true,
			err:               errors.New(""),
			response: response{
				code: http.StatusInternalServerError,
				body: "{\"text\":\"can't delete segment\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().DeleteSegment(tc.segmentName).Return(tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/segments", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.deleteSegment(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}

func TestHttpServer_createUser(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		request           string
		userId            int64
		wantToCallStorage bool
		err               error
		response          response
	}{
		{
			name:              "OK",
			request:           "{\"id\":1}",
			userId:            1,
			wantToCallStorage: true,
			err:               nil,
			response: response{
				code: http.StatusCreated,
				body: "{\"id\":1}\n",
			},
		},
		{
			name:              "Invalid json",
			request:           "{",
			userId:            1,
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"can't unmarshal json\"}\n",
			},
		},
		{
			name:              "internal error",
			request:           "{\"id\":1}",
			userId:            1,
			wantToCallStorage: true,
			err:               errors.New(""),
			response: response{
				code: http.StatusInternalServerError,
				body: "{\"text\":\"can't create user\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().CreateUser(tc.userId).Return(tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.createUser(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}

func TestHttpServer_addSegmentsToUser(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		request           string
		user              User
		wantToCallStorage bool
		err               error
		response          response
	}{
		{
			name:    "OK",
			request: "{\"id\":1,\"segments\":[\"a\",\"b\"]}",
			user: User{
				Id:       1,
				Segments: []string{"a", "b"},
			},
			wantToCallStorage: true,
			err:               nil,
			response: response{
				code: http.StatusOK,
				body: "{\"id\":1,\"segments\":[\"a\",\"b\"]}\n",
			},
		},
		{
			name:              "Invalid json",
			request:           "{",
			user:              User{},
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"can't unmarshal json\"}\n",
			},
		},
		{
			name:    "internal error",
			request: "{\"id\":1,\"segments\":[\"a\",\"b\"]}",
			user: User{
				Id:       1,
				Segments: []string{"a", "b"},
			},
			wantToCallStorage: true,
			err:               errors.New(""),
			response: response{
				code: http.StatusInternalServerError,
				body: "{\"text\":\"can't add segments to user\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().AddSegmentsToUser(tc.user).Return(tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/users", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.addSegmentsToUser(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}

func TestHttpServer_deleteSegmentsFromUser(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		request           string
		user              User
		wantToCallStorage bool
		err               error
		response          response
	}{
		{
			name:    "OK",
			request: "{\"id\":1,\"segments\":[\"a\",\"b\"]}",
			user: User{
				Id:       1,
				Segments: []string{"a", "b"},
			},
			wantToCallStorage: true,
			err:               nil,
			response: response{
				code: http.StatusOK,
				body: "{\"id\":1,\"segments\":[\"a\",\"b\"]}\n",
			},
		},
		{
			name:              "Invalid json",
			request:           "{",
			user:              User{},
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"can't unmarshal json\"}\n",
			},
		},
		{
			name:    "internal error",
			request: "{\"id\":1,\"segments\":[\"a\",\"b\"]}",
			user: User{
				Id:       1,
				Segments: []string{"a", "b"},
			},
			wantToCallStorage: true,
			err:               errors.New(""),
			response: response{
				code: http.StatusInternalServerError,
				body: "{\"text\":\"can't delete segments from user\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().DeleteSegmentsFromUser(tc.user).Return(tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/users", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.deleteSegmentsFromUser(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}

// Need to create test's with real server to check Get request and routing
func TestHttpServer_getUser(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		url               string
		userId            int64
		wantToCallStorage bool
		err               error
		user              *User
		response          response
	}{
		{
			name:              "Invalid query",
			url:               "/users/abcd",
			userId:            1,
			wantToCallStorage: false,
			err:               nil,
			user:              nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"invalid user id\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().GetUser(tc.userId).Return(tc.user, tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.getUser(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}

// Need to create test's with real server to check Get request and routing
func TestHttpServer_userHistory(t *testing.T) {
	type response struct {
		code int
		body string
	}
	tt := []struct {
		name              string
		url               string
		id                int64
		wantToCallStorage bool
		err               error
		response          response
	}{
		{
			name:              "Invalid query",
			url:               "/users/a/history",
			id:                0,
			wantToCallStorage: false,
			err:               nil,
			response: response{
				code: http.StatusBadRequest,
				body: "{\"text\":\"invalid user id\"}\n",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := NewMockStorage(ctrl)
			if tc.wantToCallStorage {
				storage.EXPECT().GetUserHistory(tc.id).Return(nil, tc.err)
			}

			server := &HttpServer{
				storage: storage,
				logger:  zap.NewNop().Sugar(),
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, tc.url, nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := server.userHistory(c)
			require.NoError(t, err)
			assert.Equal(t, tc.response.code, rec.Code)
			assert.Equal(t, tc.response.body, rec.Body.String())
		})
	}
}
