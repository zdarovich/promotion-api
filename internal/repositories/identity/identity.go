package identity

import (
	"fmt"
	"github.com/zdarovich/promotion-api/internal/cache/redis"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/service/identity"
)

type (
	// Identity struct
	Identity struct {
		Configuration   *config.Configuration
		IdentityService identity.IIdentity
	}
	// IIdentity interaface
	IIdentity interface {
		GetJWT(clientCode string, sessionKey string) (string, error)
	}
)

// New return configured identity
func New(configuration *config.Configuration) IIdentity {

	return &Identity{
		Configuration:   configuration,
		IdentityService: identity.New(configuration),
	}
}

// GetJWT returns the identity JWT token
func (identity *Identity) GetJWT(clientCode string, sessionKey string) (string, error) {

	r := redis.New(identity.Configuration)
	key := fmt.Sprintf("identity_token_%s_%s", clientCode, sessionKey)
	var res interface{} = ""

	exists, _ := r.Exists(key)
	if exists == 1 {
		res, err := r.Get(key)

		if err != nil {
			return res.(string), err
		}
	} else {
		// Request new JWT from identity integration.
		// Save it to redis and return.
		jwt, err := identity.IdentityService.GenerateNewJWT(clientCode, sessionKey)

		if err != nil {
			return res.(string), err
		}

		err = r.Set(key, jwt)

		if err != nil {
			return res.(string), err
		}

		return jwt, nil
	}

	return res.(string), nil
}
