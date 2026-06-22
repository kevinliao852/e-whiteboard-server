package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/authstate"
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

type WhiteboardSummaryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateWhiteboardRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (wc *WhiteboardController) GetWhiteboardByUserId(c *gin.Context) {
	identity, ok := authstate.FromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if identity.IsGuest {
		c.JSON(http.StatusOK, []WhiteboardSummaryResponse{})
		return
	}

	whiteboards, err := wc.service.GetUserWhiteboards(uint(identity.UserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create whiteboard"})
		return
	}

	response := make([]WhiteboardSummaryResponse, 0, len(whiteboards))
	for _, wb := range whiteboards {
		response = append(response, WhiteboardSummaryResponse{
			ID:        wb.Id,
			Name:      wb.Name,
			CreatedAt: wb.CreatedAt,
			UpdatedAt: wb.UpdatedAt,
		})
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

	identity, ok := authstate.FromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if identity.IsGuest {
		c.JSON(http.StatusForbidden, gin.H{"error": "guest cannot create whiteboard"})
		return
	}

	whiteboard, err := wc.service.CreateWhiteboard(
		core.Whiteboard{
			UserId:    uint(identity.UserID),
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
	identity, ok := authstate.FromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if identity.IsGuest {
		c.JSON(http.StatusForbidden, gin.H{"error": "guest cannot delete whiteboard"})
		return
	}

	whiteboardID, err := strconv.Atoi(c.Param("id"))
	if err != nil || whiteboardID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid whiteboard id"})
		return
	}

	if err := wc.service.DeleteWhiteboard(uint(whiteboardID)); err != nil {
		c.String(http.StatusBadRequest, "Delete whiteboard service failed")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Whiteboard deleted successfully"})
}
