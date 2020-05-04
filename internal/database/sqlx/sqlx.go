package sqlx

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Mysql driver
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/log"
)

type (
	// Mysql struct
	Mysql struct {
		Configuration *config.Configuration
		DB            *sqlx.DB
	}

	// IDB interface
	IDB interface {
		Beginx() (*sqlx.Tx, error)
		Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
		QueryRowx(query string, args ...interface{}) (*sqlx.Row, error)
		NamedExec(query string, arg interface{}) (sql.Result, error)
		Close() error
	}
)

// New returns configured mysql struct
func New(configuration *config.Configuration) IDB {

	return &Mysql{
		Configuration: configuration,
	}
}

// Close ...
func (mysql *Mysql) Close() error {
	return mysql.DB.Close()
}

// Beginx ...
func (mysql *Mysql) Beginx() (*sqlx.Tx, error) {
	err := mysql.Connect()

	if err != nil {
		mysql.logError(err.Error(), "BEGIN")
		return nil, errors.New(errorcodes.CodeDatabase)
	}

	return mysql.DB.Beginx()
}

// Queryx ...
func (mysql *Mysql) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	err := mysql.Connect()

	if err != nil {
		mysql.logError(err.Error(), query)
		return nil, errors.New(errorcodes.CodeDatabase)
	}

	defer mysql.DB.Close()
	return mysql.DB.Queryx(query, args...)
}

// QueryRowx ...
func (mysql *Mysql) QueryRowx(query string, args ...interface{}) (*sqlx.Row, error) {
	err := mysql.Connect()

	if err != nil {
		mysql.logError(err.Error(), query)
		return nil, errors.New(errorcodes.CodeDatabase)
	}

	defer mysql.DB.Close()
	return mysql.DB.QueryRowx(query, args...), nil
}

// NamedExec query returns database rows result when successful
func (mysql *Mysql) NamedExec(query string, args interface{}) (sql.Result, error) {
	err := mysql.Connect()

	if err != nil {
		mysql.logError(err.Error(), query)
		return nil, errors.New(errorcodes.CodeDatabase)
	}

	defer mysql.DB.Close()
	return mysql.DB.NamedExec(query, args)
}

// Connect opens a connection to the configured database
func (mysql *Mysql) Connect() error {

	db, err := sqlx.Open(
		mysql.Configuration.Database.Driver,
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
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
	db.Mapper = reflectx.NewMapper("json")
	mysql.DB = db
	return nil
}

// logs errors
func (mysql *Mysql) logError(errorMessage interface{}, query string) {
	log.Error(errorMessage)
}
