package identity

import (
	"github.com/zdarovich/promotion-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {

	id := New(&config.Configuration{})
	assert.NotNil(t, id)
}
