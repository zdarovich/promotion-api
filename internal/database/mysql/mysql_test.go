package mysql

import (
	"database/sql"
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDB struct{}

var queryCalled bool = false
var failQuery bool = false

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	queryCalled = true
	if failQuery {
		return nil, errors.New(errorcodes.CodeDatabase)
	}
	return &sql.Rows{}, nil
}

var queryRowCalled bool = false

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	queryRowCalled = true
	return &sql.Row{}
}

var failClose bool = false

func (m *MockDB) Close() error {
	if failClose {
		return errors.New(errorcodes.CodeDatabase)
	}
	return nil
}

/*

	Most functions in this package wrap sql functions.
	No need to cover with tests.

*/
func Test_New(t *testing.T) {

	configuration := &config.Configuration{}

	mysql := New(configuration)

	assert.NotNil(t, mysql)
}
