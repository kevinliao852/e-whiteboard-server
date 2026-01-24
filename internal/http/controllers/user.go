package controllers

import (
	"net/http"

	"github.com/kevinliao852/e-whiteboard-server/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
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
