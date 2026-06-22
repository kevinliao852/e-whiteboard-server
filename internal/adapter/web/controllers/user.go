package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type UserController struct {
	service core.UserService
}

func NewUserController(svc core.UserService) *UserController {
	return &UserController{
		service: svc,
	}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := ctrl.service.GetUser(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := sessionUserID(session.Get("user_id"))
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := ctrl.service.GetUser(strconv.Itoa(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"display-name": user.DisplayName,
	})
}
