package config

import (
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/database/sqlx"
)

type (
	// Repository struct
	Repository struct {
		Configuration *config.Configuration
		Database      sqlx.IDB
	}
	// IRepository interface
	IRepository interface {
		GetConfigByName(name string) (Conf, error)
	}
	// Conf structure of the user
	Conf struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Desc      string `json:"desc"`
		Value     string `json:"value"`
		Added     int    `json:"added"`
		Addedby   string `json:"addedby"`
		Changed   int    `json:"changed"`
		Changedby string `json:"changedby"`
	}
)

// New returns new configured user repository
func New(configuration *config.Configuration) IRepository {

	return &Repository{
		Configuration: configuration,
		Database:      sqlx.New(configuration),
	}
}

// GetConfigByName returns user by session key
func (repository *Repository) GetConfigByName(name string) (Conf, error) {

	query := "SELECT * FROM conf WHERE name=?"

	result, err := repository.Database.Queryx(query, name)
	var conf Conf
	if err != nil {
		return conf, err
	}
	if result.Next() {
		err = result.StructScan(&conf)
	}

	if err != nil {
		return conf, err
	}

	return conf, nil
}
