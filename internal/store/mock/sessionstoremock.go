package mock

import (
	"todo-go/internal/store/dbstore"

	"github.com/stretchr/testify/mock"
)

type SessionStoreMock struct {
	mock.Mock
}

func (m *SessionStoreMock) CreateSession(session *dbstore.Session) (*dbstore.Session, error) {
	args := m.Called(session)
	return args.Get(0).(*dbstore.Session), args.Error(1)
}

func (m *SessionStoreMock) GetUserFromSession(sessionID string, userID string) (*dbstore.User, error) {
	args := m.Called(sessionID, userID)
	return args.Get(0).(*dbstore.User), args.Error(1)
}
