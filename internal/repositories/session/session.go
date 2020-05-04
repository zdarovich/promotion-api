package session

import (
	"database/sql"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/database/mysql"
)

type (
	// Repository struct
	Repository struct {
		Configuration *config.Configuration
		Database      mysql.IMysql
	}
	// IRepository interface
	IRepository interface {
		GetSessionByKey(sessionKey string) (Session, error)
	}
	// Session structure of the session
	Session struct {
		ID      int
		User    string
		Key     string
		Device  string
		Started sql.NullInt64
		Expires sql.NullInt64
	}
)

// New returns new configured session repository
func New(configuration *config.Configuration) IRepository {

	return &Repository{
		Configuration: configuration,
		Database:      mysql.New(configuration),
	}
}

// GetSessionByKey returns the session if it exists
func (repository *Repository) GetSessionByKey(sessionKey string) (Session, error) {

	query := "SELECT `id`, `user`, `key`, `device`, `started`, `expires` FROM `session` WHERE `key` = ?"
	result, err := repository.Database.QueryRow(query, sessionKey)

	var session Session

	if err != nil {
		return session, err
	}

	err = result.Scan(
		&session.ID,
		&session.User,
		&session.Key,
		&session.Device,
		&session.Started,
		&session.Expires,
	)

	if err != nil {
		return session, err
	}

	return session, nil
}
