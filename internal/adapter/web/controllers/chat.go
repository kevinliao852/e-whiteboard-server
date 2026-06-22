package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	log "github.com/sirupsen/logrus"
)

type ChatController struct {
	roomService core.RoomService
	chatService core.ChatService
}

type GetChatMessagesQuery struct {
	RoomID string `form:"room-id" validate:"required"`
}

func NewChatController(roomService core.RoomService, chatService core.ChatService) *ChatController {
	return &ChatController{
		roomService: roomService,
		chatService: chatService,
	}
}

func (ctrl *ChatController) GetChatMessages(c *gin.Context) {
	var query GetChatMessagesQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	if err := validate.Struct(query); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ctrl.chatService.ListMessages(query.RoomID))
}

func (ctrl *ChatController) Chat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("id")
		if roomID == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "room id is required"})
			return
		}

		if _, err := ctrl.roomService.CreateRoom(roomID); err != nil {
			log.Printf("failed to create or load chat room: %v", err)
			return
		}

		participant, err := createParticipant(ctx)
		if err != nil {
			log.Printf("failed to create new chat participant: %v", err)
			return
		}

		if err := ctrl.roomService.JoinRoom(roomID, participant); err != nil {
			log.Printf("failed to join chat room: %v", err)
			_ = participant.Close()
			return
		}

		defer func() {
			if err := ctrl.roomService.LeaveRoom(roomID, participant); err != nil {
				log.Printf("failed to leave chat room: %v", err)
			}
			_ = participant.Close()
		}()

		participant.ReadMessage(ctx.Request.Context(), func(message []byte) {
			rawMessage := string(message)
			if err := ctrl.roomService.BroadcastToRoom(roomID, rawMessage); err != nil {
				log.Printf("failed to broadcast chat message: %v", err)
			}
			ctrl.chatService.AppendMessage(roomID, rawMessage)
		})
	}
}
