package controllers

import (
	"app/models"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {

	session := sessions.Default(c)
	fmt.Println(session.Get("id"))
	var user []models.User
	err := models.GetAllUsers(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
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

func CreateAUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	err := models.CreateAUser(&user)
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
