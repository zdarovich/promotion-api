package root

import (
	"github.com/zdarovich/promotion-api/internal/api/response"
	"github.com/zdarovich/promotion-api/internal/config"
)

type (
	// IRoot interface
	IRoot interface {
		Handle(context IGinContext) (*response.Data, error)
	}
	// IGinContext gin context interface
	IGinContext interface {
		PostForm(key string) string
	}
)

// Root struct
type Root struct{}

// New return configured root
func New(configuration *config.Configuration) IRoot {

	return &Root{}
}

// Handle handels root
func (root *Root) Handle(context IGinContext) (*response.Data, error) {

	return &response.Data{}, nil
}
