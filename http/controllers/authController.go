package controllers

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	"net/http"
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

		
		// Check if database have this user's credential.
		// If not, sign up for this user

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
