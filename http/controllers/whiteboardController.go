package controllers

import (
	"app/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWhitBoardByUser(c *gin.Context) {
	var whiteboards []models.Whiteboard
	userId, err := strconv.Atoi(c.DefaultQuery("userId", ""))

	if err != nil {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	err = models.GetWhiteboardsByUserId(&whiteboards, uint(userId))

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, whiteboards)
	}

}
