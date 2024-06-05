package todos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp/syntax"
	"strconv"
	"strings"
	"testing"
	"todo-go/internal/middleware"
	"todo-go/internal/store"
	"todo-go/internal/store/dbstore"
	"todo-go/internal/store/mock"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewTodosHandler(t *testing.T) {
	todoStore := &mock.TodoStoreMock{}
	type args struct {
		TodoStore store.TodoStore
	}
	tests := []struct {
		name string
		args args
		want *TodosHandler
	}{
		{
			name: "Test NewTodosHandler",
			args: args{
				TodoStore: todoStore,
			},
			want: &TodosHandler{todoStore},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTodosHandler(tt.args.TodoStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTodosHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodosHandler_AddTodoTemplate(t *testing.T) {
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
		todoError          error
		description        string
		expectedBody       []byte
		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "Add Todo Ok",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          nil,
			description:        "this is a description",
			expectedBody:       []byte(`this is a description`),
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    nil,
		},
		{
			name:               "No user",
			ctx:                context.Background(),
			userError:          errors.New("no user"),
			todoError:          nil,
			description:        "this is a description",
			expectedBody:       nil,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"HX-Redirect": "/login"},
		},
		{
			name:               "Error create todo",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          gorm.ErrInvalidData,
			description:        "this is a description",
			expectedBody:       []byte(`There was an error`),
			expectedStatusCode: http.StatusInternalServerError,
			expectedHeaders:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			todoStore := &mock.TodoStoreMock{}
			if tt.userError == nil {
				todoStore.On("CreateTodo", &dbstore.Todo{Description: tt.description, UserID: user.ID}).Return(&dbstore.Todo{ID: 1, Description: tt.description}, tt.todoError)
			}

			h := &TodosHandler{todoStore}

			req, _ := http.NewRequestWithContext(tt.ctx, "POST", "/", bytes.NewBufferString("description="+tt.description))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rw := httptest.NewRecorder()

			h.AddTodoTemplate(rw, req)

			assert.Equal(tt.expectedStatusCode, rw.Code, "handler returned wrong status code: got %v want %v", rw.Code, tt.expectedStatusCode)
			if tt.expectedBody != nil {
				assert.True(bytes.Contains(rw.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rw.Body.String(), tt.expectedBody)
			}
			if tt.expectedHeaders != nil {
				for k, v := range tt.expectedHeaders {
					assert.Equal(v, rw.Header().Get(k), "handler returned unexpected header: got %v want %v", rw.Header().Get(k), v)
				}
			}
		})
	}
}

func TestTodosHandler_DeleteTodoTemplate(t *testing.T) {
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
		todoError          error
		id                 string
		expectedBody       []byte
		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "Delete Todo Ok",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          nil,
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    nil,
		},
		{
			name:               "No user",
			ctx:                context.Background(),
			userError:          errors.New("no user"),
			todoError:          nil,
			id:                 "1",
			expectedBody:       nil,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"HX-Redirect": "/login"},
		},
		{
			name:               "Error delete todo",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          gorm.ErrInvalidData,
			id:                 "1",
			expectedBody:       []byte(`There was an error`),
			expectedStatusCode: http.StatusInternalServerError,
			expectedHeaders:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			todoStore := &mock.TodoStoreMock{}
			if tt.userError == nil {
				todoStore.On("DeleteTodo", tt.id, strconv.FormatUint(uint64(user.ID), 10)).Return(tt.todoError)
			}

			h := &TodosHandler{todoStore}
			req, _ := http.NewRequestWithContext(tt.ctx, "POST", "/", strings.NewReader("id="+tt.id))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rw := httptest.NewRecorder()

			h.DeleteTodoTemplate(rw, req)

			assert.Equal(tt.expectedStatusCode, rw.Code, "handler returned wrong status code: got %v want %v", rw.Code, tt.expectedStatusCode)
			if tt.expectedBody != nil {
				assert.True(bytes.Contains(rw.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rw.Body.String(), tt.expectedBody)
			}
			if tt.expectedHeaders != nil {
				for k, v := range tt.expectedHeaders {
					assert.Equal(v, rw.Header().Get(k), "handler returned unexpected header: got %v want %v", rw.Header().Get(k), v)
				}
			}
		})
	}
}

func TestTodosHandler_FinishTodoTemplate(t *testing.T) {
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
		todoError          error
		id                 string
		expectedBody       []byte
		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "Finish Todo Ok",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          nil,
			id:                 "1",
			expectedBody:       []byte(`hx-patch="unfinish-todo"`),
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    nil,
		},
		{
			name:               "No user",
			ctx:                context.Background(),
			userError:          errors.New("no user"),
			todoError:          nil,
			id:                 "1",
			expectedBody:       nil,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"HX-Redirect": "/login"},
		},
		{
			name:               "Error finish todo",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          gorm.ErrInvalidData,
			id:                 "1",
			expectedBody:       []byte(`There was an error`),
			expectedStatusCode: http.StatusInternalServerError,
			expectedHeaders:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			todoStore := &mock.TodoStoreMock{}
			if tt.userError == nil {
				todoStore.On("FinishTodo", tt.id, strconv.FormatUint(uint64(user.ID), 10)).Return(&dbstore.Todo{}, tt.todoError)
			}

			h := &TodosHandler{todoStore}
			req, _ := http.NewRequestWithContext(tt.ctx, "POST", "/", strings.NewReader("id="+tt.id))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rw := httptest.NewRecorder()

			h.FinishTodoTemplate(rw, req)

			assert.Equal(tt.expectedStatusCode, rw.Code, "handler returned wrong status code: got %v want %v", rw.Code, tt.expectedStatusCode)
			if tt.expectedBody != nil {
				assert.True(bytes.Contains(rw.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rw.Body.String(), tt.expectedBody)
			}
			if tt.expectedHeaders != nil {
				for k, v := range tt.expectedHeaders {
					assert.Equal(v, rw.Header().Get(k), "handler returned unexpected header: got %v want %v", rw.Header().Get(k), v)
				}
			}
		})
	}
}

func TestTodosHandler_UnFinishTodoTemplate(t *testing.T) {
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
		todoError          error
		id                 string
		expectedBody       []byte
		expectedStatusCode int
		expectedHeaders    map[string]string
	}{
		{
			name:               "Finish Todo Ok",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          nil,
			id:                 "1",
			expectedBody:       []byte(`hx-patch="finish-todo"`),
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    nil,
		},
		{
			name:               "No user",
			ctx:                context.Background(),
			userError:          errors.New("no user"),
			todoError:          nil,
			id:                 "1",
			expectedBody:       nil,
			expectedStatusCode: http.StatusOK,
			expectedHeaders:    map[string]string{"HX-Redirect": "/login"},
		},
		{
			name:               "Error finish todo",
			ctx:                context.WithValue(context.Background(), middleware.UserKey, user),
			userError:          nil,
			todoError:          gorm.ErrInvalidData,
			id:                 "1",
			expectedBody:       []byte(`There was an error`),
			expectedStatusCode: http.StatusInternalServerError,
			expectedHeaders:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			todoStore := &mock.TodoStoreMock{}
			if tt.userError == nil {
				todoStore.On("UnFinishTodo", tt.id, strconv.FormatUint(uint64(user.ID), 10)).Return(&dbstore.Todo{}, tt.todoError)
			}

			h := &TodosHandler{todoStore}
			req, _ := http.NewRequestWithContext(tt.ctx, "POST", "/", strings.NewReader("id="+tt.id))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rw := httptest.NewRecorder()

			h.UnFinishTodoTemplate(rw, req)

			assert.Equal(tt.expectedStatusCode, rw.Code, "handler returned wrong status code: got %v want %v", rw.Code, tt.expectedStatusCode)
			if tt.expectedBody != nil {
				assert.True(bytes.Contains(rw.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rw.Body.String(), tt.expectedBody)
			}
			if tt.expectedHeaders != nil {
				for k, v := range tt.expectedHeaders {
					assert.Equal(v, rw.Header().Get(k), "handler returned unexpected header: got %v want %v", rw.Header().Get(k), v)
				}
			}
		})
	}
}

func TestTodosHandler_AddTodo(t *testing.T) {

	tests := []struct {
		name               string
		userID             string
		userError          error
		todoBody           string
		todoBodyError      error
		todoCreated        *dbstore.Todo
		todoCreatedError   error
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name:               "Add Todo Ok",
			userID:             "1",
			userError:          nil,
			todoBody:           `{"description":"test"}`,
			todoBodyError:      nil,
			todoCreated:        &dbstore.Todo{ID: 2, Description: "test", Done: false, UserID: 1},
			todoCreatedError:   nil,
			expectedBody:       []byte("{\"id\":2,\"description\":\"test\",\"done\":false}\n"),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Bad Path Variable UserID",
			userID:             "text",
			userError:          &syntax.Error{},
			todoBody:           "",
			todoBodyError:      nil,
			todoCreated:        nil,
			todoCreatedError:   nil,
			expectedBody:       []byte("Failed to parse userID"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Bad Body",
			userID:             "1",
			userError:          nil,
			todoBody:           "{{}",
			todoBodyError:      &syntax.Error{},
			todoCreated:        nil,
			todoCreatedError:   nil,
			expectedBody:       []byte("Failed to decode todo"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Bad Creation",
			userID:             "1",
			userError:          nil,
			todoBody:           `{"description":"test"}`,
			todoBodyError:      nil,
			todoCreated:        nil,
			todoCreatedError:   gorm.ErrInvalidData,
			expectedBody:       []byte("Failed to create todo"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			todoStore := &mock.TodoStoreMock{}
			if tt.userError == nil {
				todoStore.On("CreateTodo", &dbstore.Todo{UserID: 1, Description: "test"}).Return(tt.todoCreated, tt.todoCreatedError)
			}

			h := &TodosHandler{todoStore}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userId", tt.userID)

			req, _ := http.NewRequestWithContext(
				context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
				"POST",
				fmt.Sprintf("/users/%s/todos", tt.userID),
				strings.NewReader(tt.todoBody))

			req.Header.Set("Content-Type", "application/json")

			rw := httptest.NewRecorder()

			h.AddTodo(rw, req)

			assert.Equal(tt.expectedStatusCode, rw.Code, "handler returned wrong status code: got %v want %v", rw.Code, tt.expectedStatusCode)
			if tt.expectedBody != nil {
				assert.True(bytes.Contains(rw.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rw.Body.String(), string(tt.expectedBody))
			}

		})
	}
}

func TestTodosHandler_GetAllTodos(t *testing.T) {

	tests := []struct {
		name               string
		userID             string
		userError          error
		todosCreated       *[]dbstore.Todo
		todosCreatedError  error
		expectedBody       []byte
		expectedStatusCode int
	}{
		{
			name:               "Get Todos Ok",
			userID:             "1",
			userError:          nil,
			todosCreated:       &[]dbstore.Todo{},
			todosCreatedError:  nil,
			expectedBody:       []byte("[]"),
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			todoStore := &mock.TodoStoreMock{}
			if tt.todosCreatedError == nil {
				todoStore.On("GetAllTodosFromUser", tt.userID).Return(tt.todosCreated, tt.todosCreatedError)
			}

			h := &TodosHandler{todoStore}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userId", tt.userID)

			req, _ := http.NewRequestWithContext(
				context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
				"GET",
				fmt.Sprintf("/users/%s/todos", tt.userID),
				nil)


			rw := httptest.NewRecorder()

			h.GetAllTodos(rw, req)

			assert.Equal(tt.expectedStatusCode, rw.Code, "handler returned wrong status code: got %v want %v", rw.Code, tt.expectedStatusCode)
			if tt.expectedBody != nil {
				assert.True(bytes.Contains(rw.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rw.Body.String(), string(tt.expectedBody))
			}
		})
	}
}

func TestTodosHandler_GetTodo(t *testing.T) {
	type fields struct {
		todoStore store.TodoStore
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TodosHandler{
				todoStore: tt.fields.todoStore,
			}
			h.GetTodo(tt.args.w, tt.args.r)
		})
	}
}

func TestTodosHandler_DeleteTodo(t *testing.T) {
	type fields struct {
		todoStore store.TodoStore
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TodosHandler{
				todoStore: tt.fields.todoStore,
			}
			h.DeleteTodo(tt.args.w, tt.args.r)
		})
	}
}

func TestTodosHandler_UpdateTodo(t *testing.T) {
	type fields struct {
		todoStore store.TodoStore
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TodosHandler{
				todoStore: tt.fields.todoStore,
			}
			h.UpdateTodo(tt.args.w, tt.args.r)
		})
	}
}

func TestTodosHandler_GetTodosBySearch(t *testing.T) {
	type fields struct {
		todoStore store.TodoStore
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TodosHandler{
				todoStore: tt.fields.todoStore,
			}
			h.GetTodosBySearch(tt.args.w, tt.args.r)
		})
	}
}
