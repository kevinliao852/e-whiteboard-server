package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggerMiddleWare(c *gin.Context) {
	start := time.Now()

	c.Next()

	latency := time.Since(start).Microseconds()

	log.WithFields(log.Fields{
		"body_size":     c.Writer.Size(),
		"client_ip":     c.ClientIP(),
		"error_message": c.Errors.ByType(gin.ErrorTypePrivate).String(),
		"lantency":      fmt.Sprintf("%d%s", latency, "µs"),
		"method":        c.Request.Method,
		"path":          c.Request.URL.Path,
		"status":        c.Writer.Status(),
	}).Info("Handle request")
}
