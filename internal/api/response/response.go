package response

import (
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/log"
	"net/http"
	"strconv"
	"time"
)

const (
	// StatusOK response status ok
	StatusOK string = "ok"
	// StatusError response status error
	StatusError string = "error"
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
		Error(context IGinContext, err error) *ErrorResponse
	}
	// IGinContext gin context interface
	IGinContext interface {
		JSON(code int, obj interface{})
		PostForm(string) string
		Abort()
	}
	// Data data passed to response package
	Data struct {
		Total           int
		TotalInResponse int
		Records         interface{}
	}
	// SuccessResponse structure that will be given when the request succeeds
	SuccessResponse struct {
		Status  Status      `json:"status"`
		Records interface{} `json:"records"`
	}
	// ErrorResponse structure that will be given when the request fails
	ErrorResponse struct {
		Status Status `json:"status"`
	}
	// Status structure that holds request status info
	Status struct {
		Request            string  `json:"request"`
		RequestUnixTime    int64   `json:"requestUnixTime"`
		ResponseStatus     string  `json:"responseStatus"`
		ErrorCode          int     `json:"errorCode"`
		ErrorField         string  `json:"errorField,omitempty"`
		GenerationTimeNano int64   `json:"-"`
		GenerationTime     float64 `json:"generationTime"`
		RecordsTotal       int     `json:"recordsTotal"`
		RecordsInResponse  int     `json:"recordsInResponse"`
	}
)

// New get new response struct
func New(configuration *config.Configuration, request string) IResponse {

	return &Response{
		Configuration: configuration,
		Status: Status{
			Request:            request,
			RequestUnixTime:    time.Now().Unix(),
			GenerationTimeNano: time.Now().UnixNano(),
		},
	}
}

// OK sets successful response
func (response *Response) OK(context IGinContext, responseData *Data) *SuccessResponse {

	response.end(StatusOK, errorcodes.CodeOK)

	status := response.Status
	status.RecordsTotal = responseData.Total
	status.RecordsInResponse = responseData.TotalInResponse

	sR := &SuccessResponse{
		Status:  status,
		Records: responseData.Records,
	}

	context.JSON(
		http.StatusOK,
		sR,
	)
	return sR
}

// Error sets failed response
func (response *Response) Error(context IGinContext, err error) *ErrorResponse {
	response.end(StatusError, err.Error())
	switch e := err.(type) {
	case *errorcodes.CodeError:
		response.Status.ErrorField = e.ErrorField
		response.Status.ErrorCode = e.ErrorCode
		break
	case *errorcodes.GenericError:
		log.Error(e.Generic)
		response.Status.ErrorCode = e.ErrorCode
		break
	default:
		log.Error(err)
	}
	eR := &ErrorResponse{
		Status: response.Status,
	}

	context.JSON(
		http.StatusOK,
		eR,
	)
	context.Abort()
	return eR
}

// end sets last parameters to the response struct
func (response *Response) end(responseStatus string, errorCode string) {

	response.Status.ResponseStatus = responseStatus
	response.Status.ErrorCode, _ = strconv.Atoi(errorCode)

	val := float64(time.Now().UnixNano()-response.Status.GenerationTimeNano) / float64(time.Second)
	response.Status.GenerationTime = val
}
