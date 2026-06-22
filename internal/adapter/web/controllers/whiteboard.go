package controllers

import (
	"net/http"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

var validate = validator.New()

type WhiteboardController struct {
	service core.WhiteboardService
}

func NewWhiteboardController(svc core.WhiteboardService) *WhiteboardController {
	return &WhiteboardController{
		service: svc,
	}
}

type GetWhiteboardByUserQuery struct {
	UserID uint `form:"user-id" validate:"required,gt=0"`
}

type GetWhiteboardByIdResponse struct {
	IDs []uint `json:"ids"`
}

type CreateWhiteboardRequest struct {
	UserID uint   `json:"user_id" validate:"required,gt=0"`
	Name   string `json:"name" validate:"required,min=1,max=100"`
}

type DeleteWhiteboardRequest struct {
	WhiteboardID uint `json:"whiteboard_id" validate:"required,gt=0"`
}

func (wc *WhiteboardController) GetWhiteboardByUserId(c *gin.Context) {
	var query GetWhiteboardByUserQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	if err := validate.Struct(query); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ids, err := wc.service.GetUserWhiteboards(query.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create whiteboard"})
		return
	}

	var response GetWhiteboardByIdResponse
	for _, wb := range ids {
		response.IDs = append(response.IDs, wb.Id)
	}

	c.JSON(http.StatusOK, response)
}

func (wc *WhiteboardController) CreateWhiteboard(c *gin.Context) {
	var req CreateWhiteboardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	whiteboard, err := wc.service.CreateWhiteboard(
		core.Whiteboard{
			UserId:    req.UserID,
			Name:      req.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	if err != nil {
		c.String(http.StatusBadRequest, "Create whiteboard service failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, whiteboard)
}

func (wc *WhiteboardController) DeleteWhiteboard(c *gin.Context) {
	var req DeleteWhiteboardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	err := wc.service.DeleteWhiteboard(req.WhiteboardID)
	if err != nil {
		c.String(http.StatusBadRequest, "Delete whiteboard service failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Whiteboard deleted successfully"})
}
