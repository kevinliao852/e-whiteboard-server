package controllers

import (
	"github.com/gin-gonic/gin"
)

type RoomController struct {
	Count *int
}

func (rc *RoomController) GetCurrentRoomCount(ctx *gin.Context) {
	ctx.JSON(200, interface{}(struct {
		Count int `json:"count"`
	}{Count: *rc.Count}))
}
