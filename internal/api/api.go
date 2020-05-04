package api

import (
	"fmt"
	"github.com/zdarovich/promotion-api/internal/api/router"
	"github.com/zdarovich/promotion-api/internal/config"
)

type (
	// api struct
	api struct {
		APIRouter     router.IRouter
		Configuration *config.Configuration
	}
	// IAPI interface
	IAPI interface {
		Run()
	}
)

// New get new configured api
func New(configuration *config.Configuration, router router.IRouter) IAPI {

	return &api{
		APIRouter:     router,
		Configuration: configuration,
	}
}

// Run starts the api
func (api *api) Run() {

	apiRouter := api.APIRouter.GetEngine()
	apiRouter.Run(fmt.Sprintf(":%d", api.Configuration.Port))
}
