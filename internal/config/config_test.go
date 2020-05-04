package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {

	configFileName = "../../config-sample.yml"
	c := Get()
	assert.NotNil(t, c)
}

func Test_GetFailure(t *testing.T) {

	configFileName = "/does/not/exist/"
	assert.Panics(t, func() { Get() })
}
