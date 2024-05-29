package routes

import (
	"app/http/controllers"
	"app/http/middlewares"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Handler() *gin.Engine {

	r := gin.New()
	store := cookie.NewStore([]byte("secret"))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("HOST_AllOWORIGINS")}
	config.AllowHeaders = []string{os.Getenv("HOST_AllOWHEADERS")}
	config.AllowCredentials = true

	var wc controllers.WhiteboardController
	cnt := 0
	rc := controllers.RoomController{Count: &cnt}

	r.Use(sessions.Sessions("whiteboardsession", store))
	r.Use(middlewares.LoggerMiddleWare)
	r.Use(cors.New(config))

	v1 := r.Group("/v1")
	{
		v1.GET("/user/:id", controllers.GetUser, middlewares.AuthRequired)
		v1.GET("/user", controllers.GetUsers, middlewares.AuthRequired)
		v1.POST("/user", controllers.Register)
		v1.DELETE("/user/:name", controllers.DeleteAUser, middlewares.AuthRequired)
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
