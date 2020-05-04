package session

import (
	"errors"
	"github.com/zdarovich/promotion-api/internal/api/errorcodes"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/database/mysql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	MockDB   struct{}
	MockRows struct{}
	MockRow  struct{}
)

var connectCalled bool = false
var failConnect bool = false

func (m *MockDB) Connect() error {
	connectCalled = true
	if failConnect {
		return errors.New(errorcodes.CodeDatabase)
	}
	return nil
}

var queryCalled bool = false
var failQuery bool = false
var queryResult mysql.IROWS

func (m *MockDB) Query(query string, args ...interface{}) (mysql.IROWS, error) {
	queryCalled = true
	if failQuery {
		return nil, errors.New(errorcodes.CodeDatabase)
	}
	return queryResult, nil
}

var queryRowCalled bool = false
var failQueryRow bool = false
var queryRowResult mysql.IROW

func (m *MockDB) QueryRow(query string, args ...interface{}) (mysql.IROW, error) {
	queryRowCalled = true
	if failQueryRow {
		return nil, errors.New(errorcodes.CodeDatabase)
	}
	return queryRowResult, nil
}

var nextCalled bool = false

func (m *MockRows) Next() bool {
	if !nextCalled {
		nextCalled = true
		return true
	}
	return false
}

var scanCalled bool = false
var failScan bool = false

func (m *MockRows) Scan(dest ...interface{}) error {
	scanCalled = true
	if failScan {
		return errors.New("Failed")
	}
	return nil
}

var rowScanCalled bool = false
var failRowScan bool = false

func (m *MockRow) Scan(dest ...interface{}) error {
	rowScanCalled = true
	if failRowScan {
		return errors.New("Failed")
	}
	return nil
}

func Test_New(t *testing.T) {

	configuration := &config.Configuration{}

	session := New(configuration)

	assert.NotNil(t, session)
}

func Test_GetSessionByKey(t *testing.T) {

	failQueryRow = false
	queryRowCalled = false
	queryRowResult = new(MockRow)
	rowScanCalled = false
	failRowScan = false

	configuration := &config.Configuration{}

	session := &Repository{
		Configuration: configuration,
		Database:      new(MockDB),
	}

	result, err := session.GetSessionByKey("123")

	assert.True(t, queryRowCalled)
	assert.True(t, rowScanCalled)
	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func Test_GetSessionByKeyConnectionFail(t *testing.T) {

	failQueryRow = true
	queryRowCalled = false
	queryRowResult = new(MockRow)
	rowScanCalled = false
	failRowScan = false

	configuration := &config.Configuration{}

	session := &Repository{
		Configuration: configuration,
		Database:      new(MockDB),
	}

	_, err := session.GetSessionByKey("123")

	assert.True(t, queryRowCalled)
	assert.False(t, rowScanCalled)
	assert.NotNil(t, err)
}

func Test_GetSessionByKeyScanFail(t *testing.T) {

	failQueryRow = false
	queryRowCalled = false
	queryRowResult = new(MockRow)
	rowScanCalled = false
	failRowScan = true

	configuration := &config.Configuration{}

	session := &Repository{
		Configuration: configuration,
		Database:      new(MockDB),
	}

	_, err := session.GetSessionByKey("123")

	assert.True(t, queryRowCalled)
	assert.True(t, rowScanCalled)
	assert.NotNil(t, err)
}
