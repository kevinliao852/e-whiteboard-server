package controllers

import (
	"app/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WhiteboardControllerI interface {
	GetWhiteboardByUser()
}

type WhiteboardController struct {
	model models.Whiteboard
}

func (wc *WhiteboardController) GetWhiteboardByUserId(c *gin.Context) {
	var whiteboards []models.Whiteboard
	userId, err := strconv.Atoi(c.DefaultQuery("user-id", ""))

	if err != nil {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	err = wc.model.GetWhiteboardsByUserId(&whiteboards, uint(userId))

	var whiteboardResponse []map[string]interface{}

	for i := 0; i < len(whiteboards); i++ {
		whiteboardResponse = append(whiteboardResponse, map[string]interface{}{
			"id": whiteboards[i].Id,
		})
	}

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, whiteboardResponse)
	}

}

func (wc *WhiteboardController) CreateWhiteboard(c *gin.Context) {
	var whiteboard models.Whiteboard

	err := c.BindJSON(&whiteboard)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = wc.model.CreateAWhiteboard(&whiteboard)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, whiteboard)
	}
}
