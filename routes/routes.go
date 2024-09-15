package routes

import (
	"app/http/controllers"
	"app/http/middlewares"
	"os"

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
	} else {
		log.Info("CORS is not activated")
	}

	var wc controllers.WhiteboardController
	cnt := 0
	rc := controllers.RoomController{Count: &cnt}

	r.Use(sessions.Sessions("whiteboardsession", store))
	r.Use(middlewares.LoggerMiddleWare)

	v1 := r.Group("/v1")
	{
		v1.GET("/user/:id", controllers.GetUser, currentAuthMiddleware)
		v1.GET("/user", controllers.GetUsers, currentAuthMiddleware)
		v1.POST("/user", controllers.Register)
		v1.DELETE("/user/:name", controllers.DeleteAUser, currentAuthMiddleware)
		v1.GET("/whiteboards", wc.GetWhiteboardByUserId)
		v1.POST("/whiteboards", wc.CreateWhiteboard)
		v1.DELETE("/whiteboards/:id", wc.DeleteWhiteboard)
	}

	r.POST("/login", controllers.Login(os.Getenv("GOOGLE_CLIENT_ID")))
	r.GET("/test", rc.GetCurrentRoomCount)

	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/chatting", controllers.WebsocketRoute())
		wsGroup.GET("/drawing/:id", controllers.WebsocketRoute())
	}

	return r
}
