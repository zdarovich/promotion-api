package config

import (
	"database/sql"
	"testing"

	"github.com/zdarovich/promotion-api/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type (
	databaseMock struct{}
)

func (d *databaseMock) Beginx() (*sqlx.Tx, error) { return nil, nil }

var queryXq string
var queryXa []interface{}

func (d *databaseMock) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	queryXq = query
	queryXa = args
	return new(sqlx.Rows), nil
}

func (d *databaseMock) QueryRowx(query string, args ...interface{}) (*sqlx.Row, error) {
	return nil, nil
}
func (d *databaseMock) NamedExec(query string, arg interface{}) (sql.Result, error) {
	r := new(sql.Result)

	return *r, nil
}
func (d *databaseMock) Close() error { return nil }

func TestConfig_New(t *testing.T) {
	r := New(&config.Configuration{})
	assert.IsType(t, &Repository{}, r)
}
