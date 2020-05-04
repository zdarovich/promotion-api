package api

import (
	"net/http"
	"testing"

	"github.com/zdarovich/promotion-api/internal/api/router"
	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/stretchr/testify/assert"
)

type (
	MockAPIRouter struct{}
	MockAPIEngine struct{}
)

var mockAPIEngine MockAPIEngine
var engineRunCalled bool = false

func (api *MockAPIRouter) GetEngine() router.IGINEngine {
	return &mockAPIEngine
}

func (engine *MockAPIEngine) Run(addr ...string) (err error) {
	engineRunCalled = true
	return nil
}

func (engine *MockAPIEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

// Test instance creation
func Test_New(t *testing.T) {

	api := New(&config.Configuration{}, &MockAPIRouter{})
	assert.NotNil(t, api)
}

// Test api run
func Test_Run(t *testing.T) {

	mockAPIRouter := new(MockAPIRouter)
	api := api{
		APIRouter:     mockAPIRouter,
		Configuration: &config.Configuration{},
	}

	api.Run()

	assert.True(t, engineRunCalled)
}
