package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("id") == nil {
		log.Info("Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
