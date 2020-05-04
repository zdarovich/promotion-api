package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"

	_ "github.com/go-sql-driver/mysql" // Mysql driver
	"github.com/sirupsen/logrus"
)

type (
	// Mysql struct
	Mysql struct {
		Configuration *config.Configuration
		DB            IDB
	}
	// IMysql interface
	IMysql interface {
		Query(query string, args ...interface{}) (IROWS, error)
		QueryRow(query string, args ...interface{}) (IROW, error)
		Connect() error
	}
	// IDB interface
	IDB interface {
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryRow(query string, args ...interface{}) *sql.Row
		Close() error
	}
	// IROWS interface
	IROWS interface {
		Next() bool
		Scan(dest ...interface{}) error
	}
	// IROW interface
	IROW interface {
		Scan(dest ...interface{}) error
	}
)

// New returns configured mysql struct
func New(configuration *config.Configuration) IMysql {

	return &Mysql{
		Configuration: configuration,
	}
}

// Query returns database rows result when successful
func (mysql *Mysql) Query(query string, args ...interface{}) (IROWS, error) {

	err := mysql.Connect()

	if err != nil {
		mysql.logError(err.Error(), query)
		return nil, errors.New(errorcodes.CodeDatabase)
	}

	defer mysql.DB.Close()
	return mysql.DB.Query(query, args...)
}

// QueryRow returns single database row result when successful
func (mysql *Mysql) QueryRow(query string, args ...interface{}) (IROW, error) {

	err := mysql.Connect()

	if err != nil {
		mysql.logError(err.Error(), query)
		return nil, errors.New(errorcodes.CodeDatabase)
	}

	defer mysql.DB.Close()
	return mysql.DB.QueryRow(query, args...), nil
}

// Connect opens a connection to the configured database
func (mysql *Mysql) Connect() error {

	db, err := sql.Open(
		mysql.Configuration.Database.Driver,
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			mysql.Configuration.Database.Username,
			mysql.Configuration.Database.Password,
			mysql.Configuration.Database.Server,
			mysql.Configuration.Database.Port,
			mysql.Configuration.Database.Name,
		),
	)

	if err != nil {
		mysql.logError(err.Error(), "")
		return errors.New(errorcodes.CodeDatabase)
	}

	mysql.DB = db
	return nil
}

// logs errors
func (mysql *Mysql) logError(errorMessage interface{}, query string) {

	logrus.Error(errorMessage)
}
