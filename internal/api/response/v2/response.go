package response

import (
	"net/http"
	"time"

	"github.com/zdarovich/promotion-api/internal/api/errorcodes/v2"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/log"
)

type (
	// Response response
	Response struct {
		Configuration *config.Configuration
		Status        Status
	}
	// IResponse response
	IResponse interface {
		OK(context IGinContext, responseData *Data) *SuccessResponse
		Error(context IGinContext, httpCode int, err *errorcodes.CodeError) *ErrorResponse
	}
	// IGinContext gin context interface
	IGinContext interface {
		Query(key string) string
		PostForm(key string) string
		JSON(code int, obj interface{})
		GetHeader(string) string
		Abort()
	}
	// Data data passed to response package
	Data struct {
		Records interface{}
	}
	// SuccessResponse structure that will be given when the request succeeds
	SuccessResponse struct {
		Status  Status      `json:"status"`
		Records interface{} `json:"data"`
	}
	// ErrorResponse structure that will be given when the request fails
	ErrorResponse struct {
		Status Status `json:"status"`
	}
	// Status structure that holds request status info
	Status struct {
		RequestUnixTime  int64  `json:"requestUnixTime"`
		ResponseStatus   string `json:"responseStatus"`
		ErrorCode        int    `json:"errorCode"`
		ErrorField       string `json:"errorField,omitempty"`
		ErrorDescription string `json:"errorDescription,omitempty"`
	}
)

// New get new response struct
func New(configuration *config.Configuration) IResponse {

	return &Response{
		Configuration: configuration,
		Status: Status{
			RequestUnixTime: time.Now().Unix(),
		},
	}
}

// OK sets successful response
func (response *Response) OK(context IGinContext, responseData *Data) *SuccessResponse {

	response.end("ok", errorcodes.CodeOK)

	status := response.Status

	sR := &SuccessResponse{
		Status:  status,
		Records: responseData.Records,
	}

	log.Infof("%#v", sR)

	context.JSON(
		http.StatusOK,
		sR,
	)
	return sR
}

// Error sets failed response
func (response *Response) Error(context IGinContext, httpCode int, err *errorcodes.CodeError) *ErrorResponse {
	log.Error(err)
	response.Status.ResponseStatus = "error"
	response.Status.ErrorCode = err.ErrorCode
	response.Status.ErrorField = err.ErrorField
	response.Status.ErrorDescription = err.ErrorDescription

	eR := &ErrorResponse{
		Status: response.Status,
	}

	context.JSON(
		httpCode,
		eR,
	)
	context.Abort()
	return eR
}

// end sets last parameters to the response struct
func (response *Response) end(responseStatus string, errorCode int) {

	response.Status.ResponseStatus = responseStatus
	response.Status.ErrorCode = errorCode
	response.Status.ErrorField = ""
	response.Status.ErrorDescription = ""
}
