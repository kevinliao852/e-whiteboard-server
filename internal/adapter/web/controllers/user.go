package controllers

import (
	"net/http"

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
