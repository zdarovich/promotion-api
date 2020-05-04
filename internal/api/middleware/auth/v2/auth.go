package auth

import (
	"net/http"
	"time"

	"github.com/zdarovich/promotion-api/internal/api/errorcodes/v2"
	response2 "github.com/zdarovich/promotion-api/internal/api/response/v2"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/repositories/session"

	"github.com/gin-gonic/gin"
)

type (
	// Auth struct
	Auth struct {
		Configuration     *config.Configuration
		sessionRepository session.IRepository
	}
	// IAuth interface
	IAuth interface {
		Authenticate() gin.HandlerFunc
	}
)

// New returns configured auth
func New(configuration *config.Configuration) IAuth {

	return &Auth{
		Configuration:     configuration,
		sessionRepository: session.New(configuration),
	}
}

// Authenticate will authenticate the request
func (auth *Auth) Authenticate() gin.HandlerFunc {

	return func(context *gin.Context) {

		authenticated := auth.confirmAuthentication(context.GetHeader("sessionKey"))

		if !authenticated {
			response := response2.New(auth.Configuration)
			response.Error(context, http.StatusBadRequest, errorcodes.New("", errorcodes.CodeUnauthenticated))
			return
		}

		context.Next()
	}
}

// The function that will be doing the actual authentication
func (auth *Auth) confirmAuthentication(sessionKey string) bool {

	session, err := auth.sessionRepository.GetSessionByKey(sessionKey)

	if err != nil {
		return false
	}

	// Check that a session exists and when it does then check that
	// it is not expired
	if session.ID == 0 || session.Expires.Int64 < time.Now().Unix() {
		return false
	}

	return true
}
