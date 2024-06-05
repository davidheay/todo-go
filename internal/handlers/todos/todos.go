package todos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"todo-go/internal/middleware"
	"todo-go/internal/store"
	"todo-go/internal/store/dbstore"
	"todo-go/internal/templates"

	"github.com/go-chi/chi/v5"
)

type TodosHandler struct {
	todoStore store.TodoStore
}

func NewTodosHandler(TodoStore store.TodoStore) *TodosHandler {
	return &TodosHandler{
		todoStore: TodoStore,
	}
}

func (h *TodosHandler) AddTodoTemplate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*dbstore.User)
	if !ok {
		w.Header().Set("HX-Redirect", "/login")
		return
	}
	time.Sleep(1 * time.Second)
	description := r.PostFormValue("description")
	todoRequested := dbstore.Todo{Description: description, UserID: user.ID}
	todo, err := h.todoStore.CreateTodo(&todoRequested)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}

	c := templates.Todo(todo)
	err = c.Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}
}

func (h *TodosHandler) DeleteTodoTemplate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*dbstore.User)
	if !ok {
		w.Header().Set("HX-Redirect", "/login")
		return
	}
	id := r.PostFormValue("id")
	time.Sleep(1 * time.Second)
	err := h.todoStore.DeleteTodo(id, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TodosHandler) FinishTodoTemplate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*dbstore.User)
	if !ok {
		w.Header().Set("HX-Redirect", "/login")
		return
	}
	time.Sleep(1 * time.Second)

	id := r.PostFormValue("id")
	todo, err := h.todoStore.FinishTodo(id, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}

	c := templates.TodoDone(todo)
	err = c.Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}

}

func (h *TodosHandler) UnFinishTodoTemplate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserKey).(*dbstore.User)
	if !ok {
		w.Header().Set("HX-Redirect", "/login")
		return
	}
	time.Sleep(1 * time.Second)

	id := r.PostFormValue("id")
	todo, err := h.todoStore.UnFinishTodo(id, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}
	c := templates.Todo(todo)
	err = c.Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.ErrorAdd()
		c.Render(r.Context(), w)
		return
	}

}

// AddTodo godoc
//
//	@Summary		Add a Todo
//	@Description	Add a Todo
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	dbstore.Todo
//	@Failure		400	{object}	docs.HTTPError
//	@Failure		500	{object}	docs.HTTPError
//	@Router			/users/{id}/todos [get]
func (h *TodosHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	userIDUint, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 32)
	if err != nil {
		http.Error(w, "Failed to parse userID", http.StatusBadRequest)
		return
	}
	userID := uint(userIDUint)

	var todo dbstore.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Failed to decode todo", http.StatusBadRequest)
		return
	}
	todo.UserID = userID
	todoCreated, err := h.todoStore.CreateTodo(&todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(todoCreated)
	if err != nil {
		http.Error(w, "Failed to encode todo", http.StatusInternalServerError)
		return
	}

}

// GetAllTodos godoc
//
//	@Summary		Get all Todos
//	@Description	get all Todos
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int	true	"User ID"
//	@Success		200		{array}		dbstore.Todo
//	@Failure		500		{object}	docs.HTTPError
//	@Router			/users/{userId}/todos [get]
func (h *TodosHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	if strings.TrimSpace(userID) == "" {
		http.Error(w, "Failed to parse userID", http.StatusBadRequest)
		return
	}
	todos, err := h.todoStore.GetAllTodosFromUser(userID)

	if err != nil {
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
		return
	}
}

// FirstNumsWhichSumIsX return the first set of numbers which the sum is X
func FirstNumsWhichSumIsX(nums []int, x int) []int {
	sort.Ints(nums)
	left := 0
	right := len(nums) - 1
	for left < right {
		currSum := nums[left] + nums[right]
		if currSum == x {
			return nums[left : right+1]
		} else if currSum < x {
			left++
		} else {
			right--
		}
	}
	return nil
}

// GetTodo godoc
//
//	@Summary		Get a Todo
//	@Description	get a Todo
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			todoId	path		string	true	"Todo ID"
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	dbstore.Todo
//	@Failure		500		{object}	docs.HTTPError
//	@Router			/users/{userId}/todos/{todoId} [get]
func (h *TodosHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")
	userID := chi.URLParam(r, "userId")
	todos, err := h.todoStore.GetTodoFromUser(todoID, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Failed to encode todo", http.StatusInternalServerError)
		return
	}
}

// DeleteTodo godoc
//
//	@Summary		Delete a Todo
//	@Description	Delete a Todo
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			todoId	path	string	true	"Todo ID"
//	@Param			userId	path	string	true	"User ID"
//	@Success		204		"No Content"
//	@Failure		500		{object}	docs.HTTPError
//	@Router			/users/{userId}/todos/{todoId} [delete]
func (h *TodosHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoId := chi.URLParam(r, "todoId")
	userId := chi.URLParam(r, "userId")
	err := h.todoStore.DeleteTodo(todoId, userId)

	if err != nil {
		log.Default().Println(err)
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateTodo godoc
//
//	@Summary		Update a Todo
//	@Description	Update a Todo
//	@Tags			Todo
//	@Accept			json
//	@Produce		json
//	@Param			todoId	path		string	true	"Todo ID"
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	dbstore.Todo
//	@Failure		400		{object}	docs.HTTPError
//	@Router			/users/{userId}/todos/{todoId} [put]
func (h *TodosHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Get the todo ID from the URL parameters
	todoID := chi.URLParam(r, "todoId")
	userId := chi.URLParam(r, "userId")

	// Parse the request body to get the updated todo data
	var todo dbstore.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the todo in the database
	todoUpdated, err := h.todoStore.UpdateTodo(userId, todoID, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(todoUpdated)
	if err != nil {
		http.Error(w, "Failed to encode todo", http.StatusInternalServerError)
		return
	}
}

// GetTodosBySearch godoc
//
//	@Summary		Get Todos by search
//	@Description	Get Todos by search
//	@Tags			Todos
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		string	true	"User ID"
//	@Param			title	query		string	false	"Title"
//	@Success		200		{array}		dbstore.Todo
//	@Failure		400		{object}	docs.HTTPError
//	@Router			/users/{userId}/todos/search [get]
func (h *TodosHandler) GetTodosBySearch(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	title := r.URL.Query().Get("title")
	todos, err := h.todoStore.GetTodosBySearch(userID, title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(*todos) == 0 {
		http.Error(w, fmt.Errorf("todo with title %s not found for user with ID %s", title, userID).Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Failed to encode todos", http.StatusInternalServerError)
		return
	}

}
