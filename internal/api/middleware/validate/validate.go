package validate

import (
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	response2 "github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/gin-gonic/gin"
)

// Parameter names
const (
	ParameterClientCode string = "clientCode"
	ParameterSessionKey string = "sessionKey"
	ParameterRequest    string = "request"
)

// The list of parameters that are checked when requests come in
var requiredParameters = []string{
	ParameterClientCode,
	ParameterSessionKey,
	ParameterRequest,
}

type (
	// Validate struct
	Validate struct {
		Configuration *config.Configuration
	}
	// IValidate interface
	IValidate interface {
		RequiredParameters() gin.HandlerFunc
	}
)

// New returns validate struct
func New(configuration *config.Configuration) IValidate {

	return &Validate{
		Configuration: configuration,
	}
}

// RequiredParameters checks if the required parameters have been provided
// with the request
func (validate *Validate) RequiredParameters() gin.HandlerFunc {

	return func(context *gin.Context) {

		for _, param := range requiredParameters {

			if context.PostForm(param) == "" {

				response := response2.New(validate.Configuration, context.PostForm(ParameterRequest))
				response.Error(context, errorcodes.New(param, errorcodes.CodeRequiredParameterMissing))
				return
			}
		}

		context.Next()
	}
}
