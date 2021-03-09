package controllers

import (
	"app/wshub"
	"log"

	"github.com/gin-gonic/gin"
)

func WebsocketRoute() gin.HandlerFunc {
	h := wshub.NewHub()
	return func(ctx *gin.Context) {
		go wshub.HubRun(h)
		c, err := h.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer c.Close()
		defer delete(h.Clients, c)
		h.Register <- c

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Print("read:", err)
				break
			}
			//log.Printf("recv: %s", message)

			for client := range h.Clients {
				(*client).WriteMessage(mt, message)
			}
		}
	}
}
