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

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("HOST_AllOWORIGINS")}
	config.AllowHeaders = []string{os.Getenv("HOST_AllOWHEADERS")}
	config.AllowCredentials = true

	r.Use(sessions.Sessions("whiteboardsession", store))
	r.Use(cors.New(config))

	v1 := r.Group("/v1")
	{
		v1.GET("/user/:id", controllers.GetUser, middlewares.AuthRequired)
		v1.GET("/user", controllers.GetUsers, middlewares.AuthRequired)
		v1.POST("/user", controllers.Register)
		v1.DELETE("/user/:name", controllers.DeleteAUser, middlewares.AuthRequired)
	}

	r.POST("/login", controllers.Login(os.Getenv("GOOGLE_CLIENT_ID")))

	wsGroup := r.Group("/ws")
	{
		wsGroup.GET("/chatting", controllers.WebsocketRoute())
		wsGroup.GET("/drawing", controllers.WebsocketRoute())
	}

	return r
}
