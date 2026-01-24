package route

import "github.com/gin-gonic/gin"

type Option interface {
	apply(*options)
}

type options struct {
	AuthMiddleware func(*gin.Context)
	Secret         string
	UseCORS        bool
}

type authMiddleware struct {
	AuthMiddleware func(*gin.Context)
}

func (am *authMiddleware) apply(opt *options) {
	opt.AuthMiddleware = am.AuthMiddleware
}

func WithAuthMiddleware(m func(*gin.Context)) Option {
	return &authMiddleware{AuthMiddleware: m}
}

type corsOption bool

func (am corsOption) apply(opt *options) {
	opt.UseCORS = bool(am)
}

func WithCORS() corsOption {
	return corsOption(true)
}
