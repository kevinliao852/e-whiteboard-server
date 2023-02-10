package controllers

import (
	"app/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func Login(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		client := &http.Client{}
		tokenValidator, _ := idtoken.NewValidator(context.Background(), option.WithHTTPClient(client))

		token := c.PostForm("idtoken")
		if token == "" {
			c.JSON(http.StatusBadRequest, map[string]string{"status": "invalid data"})
			return
		}

		payload, _ := tokenValidator.Validate(context.Background(), token, id)

		if payload == nil {
			c.JSON(http.StatusBadRequest, map[string]string{"status": "invalid data"})
			return
		}

		sub := fmt.Sprintf("%v", payload.Claims["sub"])
		email := fmt.Sprintf("%v", payload.Claims["email"])
		name := fmt.Sprintf("%v", payload.Claims["name"])

		// Check if database have this user's credential.
		var user models.User
		models.GetUserByGoogleId(&user, sub)

		// If not, sign up for this user
		if user.Id == 0 {

			err := models.CreateAUser(&models.User{
				GoogleId:    sub,
				Email:       email,
				DisplayName: name,
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
				return
			}
		}

		if session.Get("id") == nil {
			session.Set("name", payload.Claims["given_name"])
			session.Set("email", payload.Claims["email"])
			session.Set("id", payload.Claims["aud"])
			session.Set("exp", payload.Claims["exp"])
			session.Save()
		}

		c.JSON(http.StatusOK, "ok")
	}
}
