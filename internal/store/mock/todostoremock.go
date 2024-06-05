package mock

import (
	"todo-go/internal/store/dbstore"

	"github.com/stretchr/testify/mock"
)

type TodoStoreMock struct {
	mock.Mock
}

func (m *TodoStoreMock) CreateTodo(todo *dbstore.Todo) (*dbstore.Todo, error) {
	args := m.Called(todo)
	return args.Get(0).(*dbstore.Todo), args.Error(1)
}
func (m *TodoStoreMock) FinishTodo(todoID string, userID string) (*dbstore.Todo, error) {
	args := m.Called(todoID, userID)
	return args.Get(0).(*dbstore.Todo), args.Error(1)
}
func (m *TodoStoreMock) UnFinishTodo(todoID string, userID string) (*dbstore.Todo, error) {
	args := m.Called(todoID, userID)
	return args.Get(0).(*dbstore.Todo), args.Error(1)
}
func (m *TodoStoreMock) UpdateTodo(userID string, todoID string, todo dbstore.Todo) (*dbstore.Todo, error) {
	args := m.Called(userID, todoID, todo)
	return args.Get(0).(*dbstore.Todo), args.Error(1)
}
func (m *TodoStoreMock) DeleteTodo(todoID string, userID string) error {
	args := m.Called(todoID, userID)
	return args.Error(0)
}
func (m *TodoStoreMock) GetTodoFromUser(todoID string, userID string) (*dbstore.Todo, error) {
	args := m.Called(todoID, userID)
	return args.Get(0).(*dbstore.Todo), args.Error(1)
}
func (m *TodoStoreMock) GetAllTodosFromUser(userID string) (*[]dbstore.Todo, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]dbstore.Todo), args.Error(1)
}
func (m *TodoStoreMock) GetTodosBySearch(userID string, title string) (*[]dbstore.Todo, error) {
	args := m.Called(userID, title)
	return args.Get(0).(*[]dbstore.Todo), args.Error(1)
}
