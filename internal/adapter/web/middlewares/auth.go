package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/authstate"
	log "github.com/sirupsen/logrus"
)

func AuthRequired(c *gin.Context) {
	identity, ok := authstate.FromContext(c)
	if !ok {
		session := authstate.DebugSessionSnapshot(c)
		log.WithFields(log.Fields{
			"cookie_present": c.Request.Header.Get("Cookie") != "",
			"session_userID": session.UserID,
			"session_guest":  session.IsGuest,
			"session_role":   session.Role,
			"session_email":  session.Email != "",
			"origin":         c.Request.Header.Get("Origin"),
			"host":           c.Request.Host,
		}).Info("Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	log.WithFields(log.Fields{
		"user_id":    identity.UserID,
		"is_guest":   identity.IsGuest,
		"has_cookie": c.Request.Header.Get("Cookie") != "",
	}).Debug("Authorized request")

	c.Next()
}
