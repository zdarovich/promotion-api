package validate

import (
	"net/http"

	"github.com/zdarovich/promotion-api/internal/api/errorcodes/v2"
	response2 "github.com/zdarovich/promotion-api/internal/api/response/v2"
	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/gin-gonic/gin"
)

// Parameter names
const (
	HeaderClientCode string = "clientCode"
	HeaderSessionKey string = "sessionKey"
)

// The list of parameters that are checked when requests come in
var requiredParameters = []string{
	HeaderClientCode,
	HeaderSessionKey,
}

type (
	// Validate struct
	Validate struct {
		Configuration *config.Configuration
	}
	// IValidate interface
	IValidate interface {
		RequiredHeaders() gin.HandlerFunc
	}
)

// New returns validate struct
func New(configuration *config.Configuration) IValidate {

	return &Validate{
		Configuration: configuration,
	}
}

// RequiredHeaders checks if the required parameters have been provided
// with the request
func (validate *Validate) RequiredHeaders() gin.HandlerFunc {

	return func(context *gin.Context) {

		for _, param := range requiredParameters {

			if context.GetHeader(param) == "" {

				response := response2.New(validate.Configuration)
				response.Error(context, http.StatusBadRequest, errorcodes.New(param, errorcodes.CodeRequiredParameterMissing))
				return
			}
		}

		context.Next()
	}
}
