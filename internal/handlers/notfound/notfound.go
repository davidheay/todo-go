package notfound

import (
	"net/http"
	"todo-go/internal/templates"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	c := templates.NotFound()
	err := templates.Layout(c, "Not Found").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
