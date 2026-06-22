package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type AuthController struct {
	service core.UserService
}

func NewAuthController(svc core.UserService) *AuthController {
	return &AuthController{
		service: svc,
	}
}

func (ac AuthController) Login(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		client := &http.Client{}
		tokenValidator, err := idtoken.NewValidator(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			log.WithError(err).Error("failed to create Google ID token validator")
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
			return
		}

		token := c.PostForm("idtoken")
		if token == "" {
			err := errors.New("missing idtoken form field")
			log.WithError(err).Warn("login rejected")
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusBadRequest, map[string]string{"status": "invalid data"})
			return
		}

		payload, err := tokenValidator.Validate(context.Background(), token, id)
		if err != nil {
			log.WithError(err).WithField("token_length", len(token)).Warn("Google ID token validation failed")
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusBadRequest, map[string]string{"status": "invalid data"})
			return
		}

		if payload == nil {
			err := errors.New("token validation returned nil payload")
			log.WithError(err).Warn("login rejected")
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusBadRequest, map[string]string{"status": "invalid data"})
			return
		}

		sub := fmt.Sprintf("%v", payload.Claims["sub"])
		email := fmt.Sprintf("%v", payload.Claims["email"])
		name := fmt.Sprintf("%v", payload.Claims["name"])

		// Check if database have this user's credential.
		user, err := ac.service.GetUserByGoogleId(sub)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).WithField("google_id", sub).Error("failed to load user by Google ID")
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
			return
		}

		// If not, sign up for this user.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &core.User{
				DisplayName: name,
				Email:       email,
				GoogleID:    sub,
			}
			err := ac.service.Register(newUser)

			if err != nil {
				log.WithError(err).WithField("google_id", sub).Error("failed to register user during login")
				_ = c.Error(err).SetType(gin.ErrorTypePrivate)
				c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
				return
			}

			user, err = ac.service.GetUserByGoogleId(sub)
			if err != nil {
				log.WithError(err).WithField("google_id", sub).Error("failed to reload user after registration")
				_ = c.Error(err).SetType(gin.ErrorTypePrivate)
				c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
				return
			}
		}

		session.Set("user_id", user.ID)
		session.Set("email", user.Email)
		session.Set("display_name", user.DisplayName)
		session.Set("google_id", user.GoogleID)
		session.Set("exp", payload.Claims["exp"])
		err = session.Save()

		if err != nil {
			log.WithError(err).WithField("user_id", user.ID).Error("failed to save login session")
			_ = c.Error(err).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusInternalServerError, map[string]string{"status": "auth failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"display-name": user.DisplayName,
		})
	}
}
