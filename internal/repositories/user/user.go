package user

import (
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
		GetUser(username string) (User, error)
		GetUserBySessionKey(sessionKey string) (User, error)
	}
	// User structure of the session
	User struct {
		ID          int
		OrgPerIDDat int
		GroupID     int
		Name        string
		ShortName   string
	}
)

// New returns new configured user repository
func New(configuration *config.Configuration) IRepository {

	return &Repository{
		Configuration: configuration,
		Database:      mysql.New(configuration),
	}
}

// GetUser returns the user if it exists
func (repository *Repository) GetUser(username string) (User, error) {

	query := "SELECT id, orgper_idDat, group_id, name, shortname FROM user WHERE name = ?"
	result, err := repository.Database.QueryRow(query, username)

	var user User

	if err != nil {
		return user, err
	}

	err = result.Scan(
		&user.ID,
		&user.OrgPerIDDat,
		&user.GroupID,
		&user.Name,
		&user.ShortName,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserBySessionKey returns user by session key
func (repository *Repository) GetUserBySessionKey(sessionKey string) (User, error) {

	query := "SELECT u.id, u.orgper_idDat, u.group_id, u.name, u.shortname " +
		"FROM session s " +
		"JOIN user u ON s.user = u.shortname " +
		"WHERE s.key = ?"
	result, err := repository.Database.QueryRow(query, sessionKey)

	var user User

	if err != nil {
		return user, err
	}

	err = result.Scan(
		&user.ID,
		&user.OrgPerIDDat,
		&user.GroupID,
		&user.Name,
		&user.ShortName,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
