package dbstore

import (
	"fmt"

	"gorm.io/gorm"
)

// Todo model info
// @Description Todo  information
type Todo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Description string `json:"description"`
	UserID      uint   `json:"-"`
	Done        bool   `json:"done"`
}

type TodoStore struct {
	db *gorm.DB
}

func NewTodoStore(DB *gorm.DB) *TodoStore {
	return &TodoStore{
		db: DB,
	}
}

func (s *TodoStore) CreateTodo(todo *Todo) (*Todo, error) {
	result := s.db.Create(todo)
	if result.Error != nil {
		return nil, fmt.Errorf("error creating todo: %w", result.Error)
	}
	return todo, nil
}

func (s *TodoStore) GetTodoFromUser(todoID string, userID string) (*Todo, error) {
	var todo Todo
	result := s.db.Where("id = ? AND user_id = ?", todoID, userID).First(&todo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("todo with ID %s not found for user with ID %s", todoID, userID)
		}
		return nil, fmt.Errorf("error retrieving todo: %w", result.Error)
	}
	return &todo, nil
}

func (s *TodoStore) GetAllTodosFromUser(userID string) (*[]Todo, error) {
	var todos []Todo
	result := s.db.Where("user_id = ?", userID).Find(&todos)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving todos: %w", result.Error)
	}
	return &todos, nil
}

func (s *TodoStore) DeleteTodo(todoID string, userID string) error {
	err := s.db.Where("id = ? AND user_id = ?", todoID, userID).Delete(&Todo{}).Error
	if err != nil {
		return fmt.Errorf("error deleting todo: %w", err)
	}
	return nil
}

func (s *TodoStore) UpdateTodo(userID string, todoID string, todo Todo) (*Todo, error) {
	// Update the todo in the database
	err := s.db.Exec("UPDATE todos description = ? WHERE id = ? and user_id = ?", todo.Description, todoID, userID).Error
	if err != nil {
		return nil, fmt.Errorf("error updating todo: %w", err)
	}

	// Retrieve the updated todo from the database
	updatedTodo, err := s.GetTodoFromUser(todoID, userID)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
}

func (s *TodoStore) FinishTodo(todoID string, userID string) (*Todo, error) {
	err := s.db.Exec("UPDATE todos set done = 1 WHERE id = ? and user_id = ?", todoID, userID).Error
	if err != nil {
		return nil, fmt.Errorf("error finishing todo: %w", err)
	}

	updatedTodo, err := s.GetTodoFromUser(todoID, userID)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
}
func (s *TodoStore) UnFinishTodo(todoID string, userID string) (*Todo, error) {
	err := s.db.Exec("UPDATE todos set done = 0 WHERE id = ? and user_id = ?", todoID, userID).Error
	if err != nil {
		return nil, fmt.Errorf("error unfinishing todo: %w", err)
	}

	updatedTodo, err := s.GetTodoFromUser(todoID, userID)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
}

func (s *TodoStore) GetTodosBySearch(userID string, title string) (*[]Todo, error) {
	var todos []Todo
	result := s.db.Raw("select id, title , description from todos where user_id= ? and title like ?", userID, "%"+title+"%").Scan(&todos)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no todos with title %q found for user with ID %s", title, userID)
		}
		return nil, fmt.Errorf("error retrieving todos: %w", result.Error)
	}
	return &todos, nil
}
