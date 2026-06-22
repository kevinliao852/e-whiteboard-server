package authstate

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Identity struct {
	UserID      int
	DisplayName string
	Email       string
	IsGuest     bool
}

type SessionSnapshot struct {
	UserID  int
	IsGuest bool
	Role    string
	Email   string
}

func FromContext(c *gin.Context) (Identity, bool) {
	return FromSession(sessions.Default(c))
}

func FromSession(session sessions.Session) (Identity, bool) {
	identity := Identity{}

	if guestValue, ok := session.Get("is_guest").(bool); ok && guestValue {
		identity.IsGuest = true
		if displayName, ok := session.Get("display_name").(string); ok && displayName != "" {
			identity.DisplayName = displayName
		} else {
			identity.DisplayName = "Guest"
		}
		return identity, true
	}

	switch v := session.Get("user_id").(type) {
	case int:
		identity.UserID = v
	case int64:
		identity.UserID = int(v)
	case uint:
		identity.UserID = int(v)
	case uint64:
		identity.UserID = int(v)
	case float64:
		identity.UserID = int(v)
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return Identity{}, false
		}
		identity.UserID = id
	default:
		return Identity{}, false
	}

	if identity.UserID <= 0 {
		return Identity{}, false
	}

	if displayName, ok := session.Get("display_name").(string); ok {
		identity.DisplayName = displayName
	}
	if email, ok := session.Get("email").(string); ok {
		identity.Email = email
	}

	return identity, true
}

func DebugSessionSnapshot(c *gin.Context) SessionSnapshot {
	session := sessions.Default(c)
	snapshot := SessionSnapshot{}

	switch v := session.Get("user_id").(type) {
	case int:
		snapshot.UserID = v
	case int64:
		snapshot.UserID = int(v)
	case uint:
		snapshot.UserID = int(v)
	case uint64:
		snapshot.UserID = int(v)
	case float64:
		snapshot.UserID = int(v)
	case string:
		if id, err := strconv.Atoi(v); err == nil {
			snapshot.UserID = id
		}
	}

	if guest, ok := session.Get("is_guest").(bool); ok {
		snapshot.IsGuest = guest
	}

	if role, ok := session.Get("role").(string); ok {
		snapshot.Role = role
	}

	if email, ok := session.Get("email").(string); ok {
		snapshot.Email = email
	}

	return snapshot
}
