package store

import "todo-go/internal/store/dbstore"

type UserStore interface {
	CreateUser(email string, password string) error
	GetUser(email string) (*dbstore.User, error)
	GetUserById(userId string) (*dbstore.User, error)
}

type SessionStore interface {
	CreateSession(session *dbstore.Session) (*dbstore.Session, error)
	GetUserFromSession(sessionID string, userID string) (*dbstore.User, error)
}

type TodoStore interface {
	CreateTodo(todo *dbstore.Todo) (*dbstore.Todo, error)
	FinishTodo(todoID string, userID string) (*dbstore.Todo, error)
	UnFinishTodo(todoID string, userID string) (*dbstore.Todo, error)
	UpdateTodo(userID string, todoID string, todo dbstore.Todo) (*dbstore.Todo, error)
	DeleteTodo(todoID string, userID string) error

	GetTodoFromUser(todoID string, userID string) (*dbstore.Todo, error)
	GetAllTodosFromUser(userID string) (*[]dbstore.Todo, error)
	GetTodosBySearch(userID string, title string) (*[]dbstore.Todo, error)
}
