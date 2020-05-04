package auth

import (
	"database/sql"
	"errors"
	"github.com/zdarovich/promotion-api/internal/config"
	"github.com/zdarovich/promotion-api/internal/repositories/session"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	MockGinContext        struct{}
	MockSessionRepository struct{}
)

var postFormReturn string = ""

func (m *MockGinContext) PostForm(key string) string {
	return postFormReturn
}

var nextCalled bool = false

func (m *MockGinContext) Next() {
	nextCalled = true
}

var getSessionFail bool = false
var getSessionResult session.Session

func (m *MockSessionRepository) GetSessionByKey(sessionKey string) (session.Session, error) {
	if getSessionFail {
		return getSessionResult, errors.New("Failure")
	}
	return getSessionResult, nil
}

func Test_New(t *testing.T) {

	configuration := config.Configuration{}
	auth := New(
		&configuration,
	)
	assert.NotNil(t, auth)
}

func Test_confirmAuthentication(t *testing.T) {

	getSessionFail = false
	getSessionResult = session.Session{
		ID:     1,
		User:   "test",
		Key:    "test",
		Device: "test",
		Started: sql.NullInt64{
			Int64: 100,
			Valid: true,
		},
		Expires: sql.NullInt64{
			Int64: time.Now().Unix() + 1000,
			Valid: true,
		},
	}

	auth := Auth{
		sessionRepository: new(MockSessionRepository),
	}

	res := auth.confirmAuthentication("test")

	assert.True(t, res)
}

func Test_confirmAuthenticationDatabaseFail(t *testing.T) {

	getSessionFail = true
	getSessionResult = session.Session{}

	auth := Auth{
		sessionRepository: new(MockSessionRepository),
	}

	res := auth.confirmAuthentication("test")

	assert.False(t, res)
}

func Test_confirmAuthenticationExpired(t *testing.T) {

	getSessionFail = false
	getSessionResult = session.Session{
		ID:     1,
		User:   "test",
		Key:    "test",
		Device: "test",
		Started: sql.NullInt64{
			Int64: 100,
			Valid: true,
		},
		Expires: sql.NullInt64{
			Int64: time.Now().Unix() - 1000,
			Valid: true,
		},
	}

	auth := Auth{
		sessionRepository: new(MockSessionRepository),
	}

	res := auth.confirmAuthentication("test")

	assert.False(t, res)
}
