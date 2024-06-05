package logout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogoutHandler(t *testing.T) {
	sessionCookieName := "session"
	logoutHandler := NewLogoutHandler(sessionCookieName)

	if logoutHandler.sessionCookieName != sessionCookieName {
		t.Errorf("sessionCookieName is not equal, want %v, got %v", sessionCookieName, logoutHandler.sessionCookieName)
	}
}

func TestLogoutHandler_Post(t *testing.T) {
	assert := assert.New(t)

	logoutHandler := NewLogoutHandler("session")

	req, _ := http.NewRequest("POST", "/", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resWriter := httptest.NewRecorder()

	logoutHandler.Post(resWriter, req)

	assert.Equal(http.StatusOK, resWriter.Code, "handler returned wrong status code: got %v want %v", resWriter.Code, http.StatusPermanentRedirect)
	
	cookies := resWriter.Result().Cookies()
	sessionCookie := cookies[0]

	assert.Equal("session", sessionCookie.Name, "handler returned wrong cookie name: got %v want %v", sessionCookie.Name, "session")
	assert.Equal(-1, sessionCookie.MaxAge, "handler returned wrong cookie MaxAge: got %v want %v", sessionCookie.Value, -1)

	assert.NotEmpty(resWriter.Header().Values("HX-Redirect"))
}
