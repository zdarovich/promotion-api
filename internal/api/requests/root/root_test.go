package root

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zdarovich/promotion-api/internal/config"
)

type ginContextMock struct{}

func (g *ginContextMock) PostForm(key string) string {
	return ""
}

func Test_New(t *testing.T) {

	root := New(&config.Configuration{})
	assert.NotNil(t, root)
}

func Test_Handle(t *testing.T) {

	root := Root{}

	data, err := root.Handle(&ginContextMock{})

	assert.Nil(t, err)
	assert.NotNil(t, data)
}
