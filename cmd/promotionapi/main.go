package main

import (
	_ "github.com/zdarovich/promotion-api/docs" // Needed for swagger doc linking
	"github.com/zdarovich/promotion-api/internal/api"
	"github.com/zdarovich/promotion-api/internal/api/requests/root"
	"github.com/zdarovich/promotion-api/internal/api/router"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/requests/deletecampaigns"
	"github.com/zdarovich/promotion-api/internal/requests/getcampaigns"
	"github.com/zdarovich/promotion-api/internal/requests/savecampaigns"
)

// @title Promotion API
// @version 1.0.3
// @description POC of promotion API
//
// @tag.name general
// @tag.name campaign
//
// @BasePath /api/v1/
func main() {

	configuration := config.Get()
	handlers := make(map[string]root.IRoot)
	handlers["getCampaigns"] = getcampaigns.New(&configuration)
	handlers["saveCampaigns"] = savecampaigns.New(&configuration)
	handlers["deleteCampaigns"] = deletecampaigns.New(&configuration)
	route := router.New(&configuration, handlers)
	apiEngine := api.New(&configuration, route)
	apiEngine.Run()

	apiEngine.Run()
}
