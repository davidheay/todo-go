package home

import (
	"net/http"
	"strconv"
	"todo-go/internal/middleware"
	"todo-go/internal/store"
	"todo-go/internal/store/dbstore"
	"todo-go/internal/templates"
)

type HomeHandler struct {
	TodoStore store.TodoStore
}

func NewHomeHandler(todoStore store.TodoStore) *HomeHandler {
	return &HomeHandler{TodoStore: todoStore}
}

func (h *HomeHandler) Get(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(middleware.UserKey).(*dbstore.User)

	if !ok {
		c := templates.GuestIndex()
		err := templates.Layout(c, "My website").Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

		return
	}
	todos, err := h.TodoStore.GetAllTodosFromUser(strconv.FormatUint(uint64(user.ID), 10))

	if err != nil {
		http.Error(w, "Error retrieving todos", http.StatusInternalServerError)
		return
	}

	c := templates.Index(todos)
	err = templates.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
