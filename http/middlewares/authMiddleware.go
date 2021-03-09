package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("id") == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}
