package routes

import (
	"os"

	"github.com/kevinliao852/e-whiteboard-server/internal/http/controllers"
	"github.com/kevinliao852/e-whiteboard-server/internal/http/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Option interface {
	apply(*options)
}

type options struct {
	AuthMiddleware func(*gin.Context)
	Secret         string
	UseCORS        bool
}

type authMiddleware struct {
	AuthMiddleware func(*gin.Context)
}

func (am *authMiddleware) apply(opt *options) {
	opt.AuthMiddleware = am.AuthMiddleware
}

func WithAuthMiddleware(m func(*gin.Context)) Option {
	return &authMiddleware{AuthMiddleware: m}
}

type corsOption bool

func (am corsOption) apply(opt *options) {
	opt.UseCORS = bool(am)
}

func WithCORS() corsOption {
	return corsOption(true)
}

func Handler(opts ...Option) *gin.Engine {
	options := options{}

	for _, o := range opts {
		o.apply(&options)
	}

	currentAuthMiddleware := middlewares.AuthRequired

	if options.AuthMiddleware != nil {
		currentAuthMiddleware = options.AuthMiddleware
		log.Info("Custom AuthMiddleware is set")
	}

	r := gin.New()
	store := cookie.NewStore([]byte("secret"))

	if options.UseCORS {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{os.Getenv("HOST_AllOWORIGINS")}
		config.AllowHeaders = []string{os.Getenv("HOST_AllOWHEADERS")}
		config.AllowCredentials = true

		r.Use(cors.New(config))
		log.Info("CORS is activated")
	}

	var wc controllers.WhiteboardController
	var userController controllers.UserController
	var authController controllers.AuthController

	r.Use(sessions.Sessions("whiteboardsession", store))
	r.Use(middlewares.LoggerMiddleWare)

	v1 := r.Group("/v1")

	// user routes
	v1.GET("/user/:id", userController.GetUser, currentAuthMiddleware)

	// whiteboard routes
	v1.GET("/whiteboards", wc.GetWhiteboardByUserId)
	v1.POST("/whiteboards", wc.CreateWhiteboard)
	v1.DELETE("/whiteboards/:id", wc.DeleteWhiteboard)

	// auth routes
	r.POST("/login", authController.Login(os.Getenv("GOOGLE_CLIENT_ID")))

	// WebSocket routes
	wsGroup := r.Group("/ws")
	wsGroup.GET("/chatting", controllers.WebsocketRoute())
	wsGroup.GET("/drawing/:id", controllers.WebsocketRoute())

	return r
}
