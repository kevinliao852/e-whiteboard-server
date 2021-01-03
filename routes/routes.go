package routes

import (
	"github.com/gin-gonic/gin"
	"app/controllers"
)
func Handler() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/user", controllers.GetUsers)
		v1.POST("/user", controllers.CreateAUser)
	}

	return r
}