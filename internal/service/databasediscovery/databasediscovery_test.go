package databasediscovery

import (
	"github.com/zdarovich/promotion-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {

	configuration := &config.Configuration{}

	result := New(configuration)

	assert.NotNil(t, result)
}
