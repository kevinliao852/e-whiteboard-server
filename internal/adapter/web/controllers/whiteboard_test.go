package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWhiteboardService is a mock implementation of core.WhiteboardService
type MockWhiteboardService struct {
	mock.Mock
}

func (m *MockWhiteboardService) CreateWhiteboard(wb core.Whiteboard) (*core.Whiteboard, error) {
	args := m.Called(wb)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Whiteboard), args.Error(1)
}

func (m *MockWhiteboardService) GetUserWhiteboards(userId uint) ([]*core.Whiteboard, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*core.Whiteboard), args.Error(1)
}

func (m *MockWhiteboardService) DeleteWhiteboard(whiteboardId uint) error {
	args := m.Called(whiteboardId)
	return args.Error(0)
}

func TestWhiteboardController_GetWhiteboardByUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		now := time.Now()
		expectedWbs := []*core.Whiteboard{
			{Id: 1, UserId: 1, Name: "WB1", CreatedAt: now, UpdatedAt: now},
			{Id: 2, UserId: 1, Name: "WB2", CreatedAt: now, UpdatedAt: now},
		}
		mockService.On("GetUserWhiteboards", uint(1)).Return(expectedWbs, nil)

		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions("testsession", store))
		router.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("user_id", 1)
			_ = session.Save()
			c.Next()
		})
		router.GET("/whiteboards", ctrl.GetWhiteboardByUserId)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/whiteboards", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp []WhiteboardSummaryResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, uint(1), resp[0].ID)
		assert.Equal(t, "WB1", resp[0].Name)
		assert.Equal(t, uint(2), resp[1].ID)
		assert.Equal(t, "WB2", resp[1].Name)

		mockService.AssertExpectations(t)
	})

	t.Run("UnauthorizedWithoutSessionUserID", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions("testsession", store))
		router.GET("/whiteboards", ctrl.GetWhiteboardByUserId)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/whiteboards", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		mockService.On("GetUserWhiteboards", uint(1)).Return(nil, errors.New("db error"))

		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions("testsession", store))
		router.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("user_id", 1)
			_ = session.Save()
			c.Next()
		})
		router.GET("/whiteboards", ctrl.GetWhiteboardByUserId)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/whiteboards", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWhiteboardController_CreateWhiteboard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		reqBody := CreateWhiteboardRequest{
			Name: "New Board",
		}

		// Match any core.Whiteboard with correct UserId and Name
		mockService.On("CreateWhiteboard", mock.MatchedBy(func(wb core.Whiteboard) bool {
			return wb.UserId == 1 && wb.Name == "New Board"
		})).Return(&core.Whiteboard{Id: 10, UserId: 1, Name: "New Board", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil)

		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions("testsession", store))
		router.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("user_id", 1)
			_ = session.Save()
			c.Next()
		})
		router.POST("/whiteboards", ctrl.CreateWhiteboard)

		body, _ := json.Marshal(reqBody)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/whiteboards", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions("testsession", store))
		router.POST("/whiteboards", ctrl.CreateWhiteboard)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/whiteboards", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UnauthorizedWithoutSessionUserID", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		store := cookie.NewStore([]byte("test-secret"))
		router.Use(sessions.Sessions("testsession", store))
		router.POST("/whiteboards", ctrl.CreateWhiteboard)

		body, _ := json.Marshal(CreateWhiteboardRequest{Name: "New Board"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/whiteboards", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestWhiteboardController_DeleteWhiteboard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		mockService.On("DeleteWhiteboard", uint(10)).Return(nil)

		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		router.DELETE("/whiteboards/:id", ctrl.DeleteWhiteboard)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/whiteboards/10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidPathID", func(t *testing.T) {
		mockService := new(MockWhiteboardService)
		ctrl := NewWhiteboardController(mockService)
		router := gin.Default()
		router.DELETE("/whiteboards/:id", ctrl.DeleteWhiteboard)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/whiteboards/not-a-number", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
