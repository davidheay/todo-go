package logout

import (
	"net/http"
	"time"
)

type LogoutHandler struct {
	sessionCookieName string
}

func NewLogoutHandler(SessionCookieName string) *LogoutHandler {
	return &LogoutHandler{
		sessionCookieName: SessionCookieName,
	}
}

func (h *LogoutHandler) Post(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})
	w.Header().Set("HX-Redirect", "/")
}
