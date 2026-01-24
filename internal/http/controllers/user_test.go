package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of core.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(id string) (*core.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.User), args.Error(1)
}

func (m *MockUserService) Register(user *core.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByGoogleId(gid string) (*core.User, error) {
	args := m.Called(gid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.User), args.Error(1)
}

func TestUserController_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		expectedUser := &core.User{
			ID:          1,
			DisplayName: "Test User",
			Email:       "test@example.com",
		}
		mockService.On("GetUser", "1").Return(expectedUser, nil)

		ctrl := NewUserController(mockService)
		router := gin.Default()
		router.GET("/user/:id", ctrl.GetUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockService := new(MockUserService)
		mockService.On("GetUser", "999").Return(nil, errors.New("user not found"))

		ctrl := NewUserController(mockService)
		router := gin.Default()
		router.GET("/user/:id", ctrl.GetUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/user/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
