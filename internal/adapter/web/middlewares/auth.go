package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/authstate"
	log "github.com/sirupsen/logrus"
)

func AuthRequired(c *gin.Context) {
	if _, ok := authstate.FromContext(c); !ok {
		log.Info("Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
