package route

import (
	"os"

	"github.com/kevinliao852/e-whiteboard-server/internal/http/controllers"
	"github.com/kevinliao852/e-whiteboard-server/internal/http/middlewares"
	"github.com/kevinliao852/e-whiteboard-server/internal/model"
	"github.com/kevinliao852/e-whiteboard-server/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	r.Use(sessions.Sessions("whiteboardsession", store))
	r.Use(middlewares.LoggerMiddleWare)

	if options.UseCORS {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{os.Getenv("HOST_AllOW_ORIGINS")}
		// config.AllowHeaders = []string{os.Getenv("HOST_AllOW_HEADERS")}
		config.AllowHeaders = []string{"Content-Type", "Authorization"}
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		config.AllowCredentials = true

		r.Use(cors.New(config))
		log.Info("CORS is activated")
	}

	whiteboardController := controllers.NewWhiteboardController(&service.WhiteboardSVC{
		Model: &model.Whiteboard{}})
	userController := controllers.NewUserController(&service.UserSVC{
		Model: &model.User{}})
	authController := controllers.NewAuthController(&service.UserSVC{
		Model: &model.User{}})
	drawingController := controllers.DrawingController{
		RoomService: service.NewRoomSVC(),
	}

	v1 := r.Group("/v1")

	// user routes
	v1.GET("/user/:id", currentAuthMiddleware, userController.GetUser)

	// whiteboard routes
	v1.GET("/whiteboards", whiteboardController.GetWhiteboardByUserId)
	v1.POST("/whiteboards", whiteboardController.CreateWhiteboard)
	v1.DELETE("/whiteboards/:id", whiteboardController.DeleteWhiteboard)

	// auth routes
	r.POST("/login", authController.Login(os.Getenv("GOOGLE_CLIENT_ID")))

	// WebSocket routes
	wsGroup := r.Group("/ws")
	wsGroup.GET("/drawing", drawingController.Draw())
	wsGroup.GET("/drawing/:id", drawingController.Draw())

	return r
}
