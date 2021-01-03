package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"app/models"
)

func GetUsers(c *gin.Context) {
	var user []models.User
	err := models.GetAllUsers(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
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