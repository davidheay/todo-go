package register

import (
	"net/http"
	"todo-go/internal/store"
	"todo-go/internal/templates"
)

type RegisterHandler struct {
	userStore store.UserStore
}

func NewRegisterHandler(UserStore store.UserStore) *RegisterHandler {
	return &RegisterHandler{
		UserStore,
	}
}

func (h *RegisterHandler) Get(w http.ResponseWriter, r *http.Request) {
	c := templates.RegisterPage()
	err := templates.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}

func (h *RegisterHandler) Post(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.userStore.CreateUser(email, password)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}
