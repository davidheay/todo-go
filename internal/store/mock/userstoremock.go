package mock

import (
	"todo-go/internal/store/dbstore"

	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (m *UserStoreMock) CreateUser(email string, password string) error {
	args := m.Called(email, password)

	return args.Error(0)
}

func (m *UserStoreMock) GetUser(email string) (*dbstore.User, error) {
	args := m.Called(email)
	return args.Get(0).(*dbstore.User), args.Error(1)
}

func (m *UserStoreMock) GetUserById(userId string) (*dbstore.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*dbstore.User), args.Error(1)
}