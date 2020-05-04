package router

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/api/middleware/auth"
	"github.com/zdarovich/promotion-api/internal/api/middleware/discovery"
	"github.com/zdarovich/promotion-api/internal/api/middleware/validate"
	"github.com/zdarovich/promotion-api/internal/api/requests/root"
	"github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const apiLogFilePath string = "/src/build/logs/"
const apiLogFileName string = "api.log"

type (
	// router router
	router struct {
		Middleware struct {
			Validate  validate.IValidate
			Discovery discovery.IDiscovery
			Auth      auth.IAuth
		}
		Configuration *config.Configuration
		Response      response.IResponse
		Root          root.IRoot
		Handlers      map[string]root.IRoot
	}
	// IRouter irouter
	IRouter interface {
		GetEngine() IGINEngine
	}
	// IGINEngine iginengine
	IGINEngine interface {
		Run(addr ...string) (err error)
		ServeHTTP(w http.ResponseWriter, req *http.Request)
	}
)

// New configures and returns router
func New(configuration *config.Configuration, handlers map[string]root.IRoot) IRouter {

	return &router{
		Middleware: struct {
			Validate  validate.IValidate
			Discovery discovery.IDiscovery
			Auth      auth.IAuth
		}{
			Validate:  validate.New(configuration),
			Discovery: discovery.New(configuration),
			Auth:      auth.New(configuration),
		},
		Configuration: configuration,
		Handlers:      handlers,
		Root:          root.New(configuration),
	}
}

// GetEngine configures and returns the router
//
// @title Picture API
// @version 1.0.3
// @description POC of picture API
//
// @tag.name general
// @tag.name product
//
// @BasePath /api
func (apiRouter *router) GetEngine() IGINEngine {

	apiRouter.configure()
	router := gin.Default()

	// API Root endpoint
	router.GET("/api", apiRouter.handleRoot)

	// Swagger endpoint
	router.GET("/documentation/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
	))

	// API v1 endpoints
	// All requests in the group use the authenticate middleware
	// thus all requests here need to be authenticated
	apiV1Group := router.Group("/api/v1")

	apiV1Group.Use(apiRouter.Middleware.Validate.RequiredParameters())
	apiV1Group.Use(apiRouter.Middleware.Discovery.Discover())
	apiV1Group.Use(apiRouter.Middleware.Auth.Authenticate())

	apiV1Group.POST("*any", apiRouter.handlePostRequest)

	return router
}

// configure internal function to set the path and the type
// of the logger and all other configurations
func (apiRouter *router) configure() {

	// Colored logs are not needed
	gin.DisableConsoleColor()

	// Disables debug in release mode
	if apiRouter.Configuration.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if !apiRouter.Configuration.LogsEnabled {
		return
	}

	// Set writer when file logging is enabled
	f, _ := os.Create(apiRouter.Configuration.LogFilePath + apiLogFileName)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

// startStatus sets status updates for request start
func (apiRouter *router) startStatus(request string) {

	apiRouter.Response = response.New(
		apiRouter.Configuration,
		request,
	)
}

// START - Functions that handle routes

// handlePostRequest handles all requests to the api v1 endpoint
// and redirects them to the correct handlers
func (apiRouter *router) handlePostRequest(context *gin.Context) {

	request := context.PostForm("request")
	apiRouter.startStatus(request)

	if handler, ok := apiRouter.Handlers[request]; !ok {
		apiRouter.handleUnknownRequest(context)
		return
	} else {
		data, err := handler.Handle(context)

		if err != nil {
			apiRouter.Response.Error(context, err)
			return
		}
		apiRouter.Response.OK(context, data)
	}
}

// handleRoot handles the Root response
// Root endpoint request handle cannot return an error
//
// @Summary API Root endpoint
// @Description Can be used to check if the API is up and running
// @Description It does not however attepmt to make any connections
// @ID handleRoot
// @tags general
// @Produce json
// @Success 200 {array} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @router / [get]
func (apiRouter *router) handleRoot(context *gin.Context) {

	apiRouter.startStatus("Root")
	data, _ := apiRouter.Root.Handle(context)

	apiRouter.Response.OK(context, data)
}

// handleUnknownRequest default handler when a call was made for a request
// that does not exist
func (apiRouter *router) handleUnknownRequest(context *gin.Context) {

	apiRouter.Response.Error(context, errors.New(errorcodes.CodeUnknownRequest))
}
