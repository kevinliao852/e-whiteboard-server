package route

import (
	"os"

	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/db"
	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/state"
	webcontrollers "github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/controllers"
	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/middlewares"
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

	whiteboardController := webcontrollers.NewWhiteboardController(&service.WhiteboardSVC{
		Model: &db.Whiteboard{}})
	userController := webcontrollers.NewUserController(&service.UserSVC{
		Model: &db.User{}})
	authController := webcontrollers.NewAuthController(&service.UserSVC{
		Model: &db.User{}})
	roomState := state.NewRoomState()
	roomController := webcontrollers.NewRoomController(roomState)
	drawingController := webcontrollers.DrawingController{
		RoomService:    roomState,
		DrawingService: service.NewDrawingSVC(&db.WhiteboardCanvasData{}),
	}

	v1 := r.Group("/v1")

	// user routes
	v1.GET("/user/:id", currentAuthMiddleware, userController.GetUser)
	v1.GET("/rooms", roomController.ListRooms)

	// whiteboard routes
	v1.GET("/whiteboards", currentAuthMiddleware, whiteboardController.GetWhiteboardByUserId)
	v1.POST("/whiteboards", currentAuthMiddleware, whiteboardController.CreateWhiteboard)
	v1.DELETE("/whiteboards/:id", currentAuthMiddleware, whiteboardController.DeleteWhiteboard)

	// auth routes
	r.POST("/login", authController.Login(os.Getenv("GOOGLE_CLIENT_ID")))

	// WebSocket routes
	wsGroup := r.Group("/ws")
	wsGroup.GET("/drawing", drawingController.Draw())
	wsGroup.GET("/drawing/:id", drawingController.Draw())

	return r
}
