package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type authUserServiceStub struct {
	mock.Mock
}

func (s *authUserServiceStub) GetUser(id string) (*core.User, error) {
	args := s.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.User), args.Error(1)
}

func (s *authUserServiceStub) Register(user *core.User) error {
	args := s.Called(user)
	return args.Error(0)
}

func (s *authUserServiceStub) GetUserByGoogleId(gid string) (*core.User, error) {
	args := s.Called(gid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.User), args.Error(1)
}

func TestAuthController_GuestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := NewAuthController(&authUserServiceStub{})
	router := gin.Default()
	store := cookie.NewStore([]byte("test-secret"))
	router.Use(sessions.Sessions("testsession", store))
	router.POST("/guest-login", ctrl.GuestLogin())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/guest-login", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"role\":\"guest\"")
	assert.Contains(t, w.Header().Get("Set-Cookie"), "testsession=")
}
