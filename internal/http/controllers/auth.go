package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

type AuthController struct {
	service *service.UserService
}

func (ac AuthController) Login(id string) gin.HandlerFunc {
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
		user, err := ac.service.GetUserByGoogleId(sub)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
			return
		}

		// If not, sign up for this user
		if user.ID == 0 {
			err := ac.service.Register(&core.User{
				DisplayName: name,
				Email:       email,
				GoogleID:    sub,
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
			err = session.Save()

			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
				return
			}

		}

		c.JSON(http.StatusOK, gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"display-name": user.DisplayName,
		})
	}
}
