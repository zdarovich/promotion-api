package response

import (
	"errors"
	"testing"

	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/stretchr/testify/assert"
)

type MockGinContext struct{}

var jsonCalled bool = false

func (m *MockGinContext) JSON(code int, obj interface{}) {
	jsonCalled = true
}

var abortCalled bool = false

func (m *MockGinContext) Abort() {
	abortCalled = true
}

func (m *MockGinContext) PostForm(string) string {
	return ""
}

func Test_New(t *testing.T) {

	configuration := &config.Configuration{}
	response := New(configuration, "Test")

	assert.NotNil(t, response)
}

func Test_OK(t *testing.T) {

	configuration := &config.Configuration{}
	response := New(configuration, "Test")

	jsonCalled = false

	context := new(MockGinContext)
	data := Data{
		Total:           1,
		TotalInResponse: 1,
		Records:         make([]string, 0),
	}

	successResponse := response.OK(context, &data)

	assert.NotNil(t, successResponse)
	assert.True(t, jsonCalled)
	assert.Equal(t, 0, successResponse.Status.ErrorCode)
	assert.Equal(t, 1, successResponse.Status.RecordsTotal)
	assert.Equal(t, 1, successResponse.Status.RecordsInResponse)
}

func Test_Error(t *testing.T) {

	configuration := &config.Configuration{}
	response := New(configuration, "Test")

	abortCalled = false

	context := new(MockGinContext)

	err := errors.New(errorcodes.CodeUnauthenticated)

	errorResponse := response.Error(context, err)

	assert.NotNil(t, errorResponse)
	assert.True(t, abortCalled)
	assert.Equal(t, 1051, errorResponse.Status.ErrorCode)
}
