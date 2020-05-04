package router

import (
	"net/http"

	"github.com/zdarovich/promotion-api/internal/api/errorcodes/v2"
	"github.com/zdarovich/promotion-api/internal/api/middleware/auth/v2"
	"github.com/zdarovich/promotion-api/internal/api/middleware/discovery/v2"
	"github.com/zdarovich/promotion-api/internal/api/middleware/validate/v2"
	"github.com/zdarovich/promotion-api/internal/api/response/v2"
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
		CRUDHandlers  []Route
	}
	// Route struct
	Route struct {
		Method      string
		Pattern     string
		HandlerFunc gin.HandlerFunc
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
	// IGINContext gin context
	IGINContext interface {
		Query(key string) string
		PostForm(key string) string
	}
)

// New configures and returns router
func New(configuration *config.Configuration, handlers []Route) IRouter {

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
		CRUDHandlers:  handlers,
	}
}

// GetEngine configures and returns the router
func (apiRouter *router) GetEngine() IGINEngine {

	apiRouter.configure()
	router := gin.Default()

	// API Root endpoint
	router.GET("/", apiRouter.handleRoot)

	// Swagger endpoint
	router.GET("/documentation/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
	))

	// API v1 endpoints
	// All requests in the group use the authenticate middleware
	// thus all requests here need to be authenticated
	apiV1Group := router.Group("/v1")

	apiV1Group.Use(apiRouter.Middleware.Validate.RequiredHeaders())
	apiV1Group.Use(apiRouter.Middleware.Discovery.Discover())
	apiV1Group.Use(apiRouter.Middleware.Auth.Authenticate())

	for _, route := range apiRouter.CRUDHandlers {
		apiV1Group.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}

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
}

// startStatus sets status updates for request start
func (apiRouter *router) startStatus() {

	apiRouter.Response = response.New(
		apiRouter.Configuration,
	)
}

// handleRoot handles the Root response
// Root endpoint request handle cannot return an error
func (apiRouter *router) handleRoot(context *gin.Context) {

	apiRouter.startStatus()

	apiRouter.Response.OK(context, &response.Data{})
}

// handleUnknownRequest default handler when a call was made for a request
// that does not exist
func (apiRouter *router) handleUnknownRequest(context *gin.Context) {

	apiRouter.Response.Error(context, http.StatusNotFound, errorcodes.New("", errorcodes.CodeUnknownRequest))
}
