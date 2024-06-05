package home

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"todo-go/internal/middleware"
	"todo-go/internal/store/dbstore"
	"todo-go/internal/store/mock"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewHomeHandler(t *testing.T) {

	todoStore := &mock.TodoStoreMock{}
	homeHandler := NewHomeHandler(todoStore)

	if homeHandler.TodoStore != todoStore {
		t.Errorf("todoStore is not equal, want %v, got %v", todoStore, homeHandler.TodoStore)

	}

}

func TestHomeHandler_Get(t *testing.T) {
	user := &dbstore.User{
		ID:       1,
		Email:    "mock@mail.com",
		Password: "1234",
		Todos:    []dbstore.Todo{},
	}
	tests := []struct {
		name               string
		ctx                context.Context
		user               *dbstore.User
		userError          error
		expectedText       []byte
		expectedStatusCode int
	}{
		{
			name:               "logged in",
			user:               user,
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			expectedText:       []byte(`add-todo`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "not logged in",
			user:               nil,
			ctx:                context.WithValue(context.Background(), middleware.UserKey, nil),
			userError:          nil,
			expectedText:       []byte(`A simple app to manage your tasks`),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "logged in error storage",
			user:               user,
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          gorm.ErrInvalidTransaction,
			expectedText:       []byte(`Error retrieving todos`),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			todoStoreMock := &mock.TodoStoreMock{}
			if tt.user != nil {
				todoStoreMock.On("GetAllTodosFromUser", strconv.FormatUint(uint64(tt.user.ID), 10)).Return(&tt.user.Todos, tt.userError)
			}

			h := &HomeHandler{
				TodoStore: todoStoreMock,
			}
			req, _ := http.NewRequestWithContext(tt.ctx, "GET", "/", nil)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resWriter := httptest.NewRecorder()

			h.Get(resWriter, req)
			assert.Equal(tt.expectedStatusCode, resWriter.Code, "handler returned wrong status code: got %v want %v", resWriter.Code, tt.expectedStatusCode)
			assert.True(bytes.Contains(resWriter.Body.Bytes(), tt.expectedText), "handler returned unexpected body: got %v want %v", resWriter.Body.String(), tt.expectedText)

			todoStoreMock.AssertExpectations(t)
		})
	}
}
