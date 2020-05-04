package identity

import (
	"encoding/json"
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/log"
	"github.com/zdarovich/promotion-api/internal/repositories/user"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	// Identity struct
	Identity struct {
		UserRepository user.IRepository
		Configuration  *config.Configuration
	}
	// IIdentity interface
	IIdentity interface {
		GenerateNewJWT(clientCode string, sessionKey string) (string, error)
	}
	identityResponse struct {
		Result struct {
			JWT string `json:"jwt"`
		} `json:"result"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
)

// New return new configured identity
func New(configuration *config.Configuration) IIdentity {

	return &Identity{
		UserRepository: user.New(configuration),
		Configuration:  configuration,
	}
}

// GenerateNewJWT generates a new JWT
func (identity *Identity) GenerateNewJWT(clientCode string, sessionKey string) (string, error) {

	var httpClient = &http.Client{
		Timeout: time.Duration(identity.Configuration.Identity.Timeout) * time.Second,
	}

	data, err := identity.getURLValues(clientCode, sessionKey)

	if err != nil {
		identity.logError("identity session key failure", data, err.Error())
		return "", errors.New(errorcodes.CodeIdentity)
	}

	url := identity.Configuration.Identity.Server

	req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	result, err := httpClient.Do(req)

	if err != nil {
		identity.logError(url, data, err.Error())
		return "", errors.New(errorcodes.CodeIdentity)
	}

	if result.StatusCode != 201 && result.StatusCode != 200 {
		identity.logError(url, data, result.StatusCode)
		return "", errors.New(errorcodes.CodeIdentity)
	}

	defer result.Body.Close()

	var iR identityResponse

	err = json.NewDecoder(result.Body).Decode(&iR)

	if err != nil {
		identity.logError(url, data, err.Error())
		return "", errors.New(errorcodes.CodeIdentity)
	}

	return iR.Result.JWT, nil
}

// composes the parameters for the token request
func (identity *Identity) getURLValues(clientCode string, sessionKey string) (url.Values, error) {

	v := url.Values{}
	user, err := identity.UserRepository.GetUserBySessionKey(sessionKey)

	if err != nil {
		return v, err
	}

	v.Set("api[jwt]", identity.Configuration.Identity.Token)
	v.Set("parameters[username]", user.ShortName)
	v.Set("parameters[clientCode]", clientCode)

	return v, nil
}

// logs the error message
func (identity *Identity) logError(url string, data url.Values, errorMessage interface{}) {

	log.Error(errorMessage)
}
