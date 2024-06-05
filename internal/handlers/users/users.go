package users

import (
	"encoding/json"
	"net/http"
	"todo-go/internal/store"
	"todo-go/internal/store/dbstore"

	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	userStore store.UserStore
}

func NewUsersHandler(UserStore store.UserStore) *UsersHandler {
	return &UsersHandler{
		userStore: UserStore,
	}
}

// GetUserById godoc
//
//	@Summary		GetUserById
//	@Description	Get User By Id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	dbstore.User
//	@Failure		400	{object}	docs.HTTPError
//	@Failure		500	{object}	docs.HTTPError
//	@Router			/users/{id} [get]
func (h *UsersHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	var user *dbstore.User
	var err error
	user, err = h.userStore.GetUserById(userID)

	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}
