package discovery

import (
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	response2 "github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/service/databasediscovery"

	"github.com/gin-gonic/gin"
)

type (
	// Discovery struct
	Discovery struct {
		Configuration     *config.Configuration
		DatabaseDiscovery databasediscovery.IDatabaseDiscovery
	}
	// IDiscovery interface
	IDiscovery interface {
		Discover() gin.HandlerFunc
	}
)

// New returns configured auth
func New(configuration *config.Configuration) IDiscovery {

	return &Discovery{
		Configuration:     configuration,
		DatabaseDiscovery: databasediscovery.New(configuration),
	}
}

// Discover based on the incoming clients code sets the database parameters to
// the configuration
func (discovery *Discovery) Discover() gin.HandlerFunc {

	return func(context *gin.Context) {

		clientCode := context.PostForm("clientCode")
		err := discovery.setDatabaseConfig(clientCode)

		if err != nil {
			response := response2.New(discovery.Configuration, context.PostForm("request"))
			response.Error(context, errors.New(errorcodes.CodeDatabase))
			return
		}

		context.Next()
	}
}

// Depending on the configuration sets the database configuration
func (discovery *Discovery) setDatabaseConfig(clientCode string) error {

	if discovery.Configuration.Database.Discovery.Enabled == false {
		// At this state the configuration has already been mapped from the
		// the default set - no need to re-set the values
		return nil
	}

	database, err := discovery.DatabaseDiscovery.GetDatabase(clientCode)

	if err != nil {
		return err
	}

	discovery.Configuration.Database.Name = database.DatabaseName
	discovery.Configuration.Database.Server = database.Host
	discovery.Configuration.Database.Port = database.Port
	discovery.Configuration.Database.Username = database.User
	discovery.Configuration.Database.Password = database.Password

	return nil
}
