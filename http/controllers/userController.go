package controllers

import (
	"app/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {

	var users []models.User
	err := models.GetAllUsers(&users)
	c.JSON(http.StatusOK, []interface{}{})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {

		c.JSON(http.StatusOK, users)
	}

}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := models.GetUserById(&user, id); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func Register(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = models.CreateAUser(&user)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func DeleteAUser(c *gin.Context) {
	if err := models.DeleteAUser(c.Param("name")); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, "success")
	}
}
