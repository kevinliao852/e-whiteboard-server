package controllers

import (
	"net/http"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"

	"github.com/go-playground/validator/v10"

	"github.com/gin-contrib/sessions"
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

type GetWhiteboardByIdResponse struct {
	IDs []uint `json:"ids"`
}

type CreateWhiteboardRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type DeleteWhiteboardRequest struct {
	WhiteboardID uint `json:"whiteboard_id" validate:"required,gt=0"`
}

func (wc *WhiteboardController) GetWhiteboardByUserId(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := sessionUserID(session.Get("user_id"))
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ids, err := wc.service.GetUserWhiteboards(uint(userID))
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

	session := sessions.Default(c)
	userID, ok := sessionUserID(session.Get("user_id"))
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	whiteboard, err := wc.service.CreateWhiteboard(
		core.Whiteboard{
			UserId:    uint(userID),
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
