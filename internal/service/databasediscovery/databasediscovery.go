package databasediscovery

import (
	"encoding/json"
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/log"
	"net/http"
	"time"
)

type (
	// DatabaseDiscovery struct
	DatabaseDiscovery struct {
		Configuration *config.Configuration
	}
	// IDatabaseDiscovery interface
	IDatabaseDiscovery interface {
		GetDatabase(clientCode string) (Database, error)
	}
	// Database dtabase struct object
	Database struct {
		Tenant       string `json:"tenant"`
		DatabaseName string `json:"databaseName"`
		Host         string `json:"host"`
		Port         int    `json:"port"`
		User         string `json:"username"`
		Password     string `json:"password"`
	}
)

// New returns new configured database discovery struct
func New(configuration *config.Configuration) IDatabaseDiscovery {

	return &DatabaseDiscovery{
		Configuration: configuration,
	}
}

// GetDatabase gets details of the database for current client
func (databasediscovery *DatabaseDiscovery) GetDatabase(clientCode string) (Database, error) {

	var httpClient = &http.Client{
		Timeout: time.Duration(databasediscovery.Configuration.Database.Discovery.Timeout) * time.Second,
	}

	result, err := httpClient.Get(
		databasediscovery.Configuration.Database.Discovery.Server + "/v1/tenants/" + clientCode,
	)

	if err != nil {
		databasediscovery.logError(err.Error())
		return Database{}, errors.New(errorcodes.CodeDBDiscovery)
	}
	defer result.Body.Close()

	var database Database
	err = json.NewDecoder(result.Body).Decode(&database)

	if err != nil {
		databasediscovery.logError(err.Error())
		return Database{}, errors.New(errorcodes.CodeDBDiscovery)
	}

	return database, nil
}

// logs errors
func (databasediscovery *DatabaseDiscovery) logError(errorMessage interface{}) {
	log.Error(errorMessage)
}
