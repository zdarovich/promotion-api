package discovery

import (
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/service/databasediscovery"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDatabaseDiscovery struct{}

var failGetDatabase bool = false
var getDatabaseResult databasediscovery.Database

func (m *MockDatabaseDiscovery) GetDatabase(clientCode string) (databasediscovery.Database, error) {
	if failGetDatabase {
		return getDatabaseResult, errors.New(errorcodes.CodeDatabase)
	}
	return getDatabaseResult, nil
}

func Test_New(t *testing.T) {

	configuration := config.Configuration{}
	result := New(&configuration)

	assert.NotNil(t, result)
}

func Test_setDatabaseConfig(t *testing.T) {

	failGetDatabase = false
	getDatabaseResult = databasediscovery.Database{
		Tenant:       "100",
		DatabaseName: "crmx_100",
		Host:         "0.0.0.0",
		Port:         3306,
	}

	configuration := config.Configuration{}
	configuration.Database.Discovery.Enabled = true
	configuration.Database.Name = "crmx_500"
	configuration.Database.Server = "1.1.1.1"
	configuration.Database.Port = 1000

	discovery := &Discovery{
		Configuration:     &configuration,
		DatabaseDiscovery: new(MockDatabaseDiscovery),
	}

	err := discovery.setDatabaseConfig("123")

	assert.Nil(t, err)
	assert.Equal(t, getDatabaseResult.DatabaseName, configuration.Database.Name)
	assert.Equal(t, getDatabaseResult.Host, configuration.Database.Server)
	assert.Equal(t, getDatabaseResult.Port, configuration.Database.Port)
}

func Test_setDatabaseConfigDiscoveryDisabled(t *testing.T) {

	failGetDatabase = false
	getDatabaseResult = databasediscovery.Database{
		Tenant:       "100",
		DatabaseName: "crmx_100",
		Host:         "0.0.0.0",
		Port:         3306,
	}

	configuration := config.Configuration{}
	configuration.Database.Discovery.Enabled = false
	configuration.Database.Name = "crmx_500"
	configuration.Database.Server = "1.1.1.1"
	configuration.Database.Port = 1000

	discovery := &Discovery{
		Configuration:     &configuration,
		DatabaseDiscovery: new(MockDatabaseDiscovery),
	}

	err := discovery.setDatabaseConfig("123")

	assert.Nil(t, err)
	assert.Equal(t, "crmx_500", configuration.Database.Name)
	assert.Equal(t, "1.1.1.1", configuration.Database.Server)
	assert.Equal(t, 1000, configuration.Database.Port)
}

func Test_setDatabaseConfigDatabaseFail(t *testing.T) {

	failGetDatabase = true
	getDatabaseResult = databasediscovery.Database{}

	configuration := config.Configuration{}
	configuration.Database.Discovery.Enabled = true

	discovery := &Discovery{
		Configuration:     &configuration,
		DatabaseDiscovery: new(MockDatabaseDiscovery),
	}

	err := discovery.setDatabaseConfig("123")

	assert.NotNil(t, err)
}
