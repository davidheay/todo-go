package notfound

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundHandler(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewNotFoundHandler())
}

func TestNotFoundHandler_NotFound(t *testing.T) {
	tests := []struct {
		name               string
		expectedText       []byte
		expectedStatusCode int
	}{
		{
			name:               "default",
			expectedText:       []byte(`Something's missing`),
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			h := &NotFoundHandler{}
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resWriter := httptest.NewRecorder()

			h.NotFound(resWriter, req)
			assert.Equal(tt.expectedStatusCode, resWriter.Code, "handler returned wrong status code: got %v want %v", resWriter.Code, tt.expectedStatusCode)
			assert.True(bytes.Contains(resWriter.Body.Bytes(), tt.expectedText), "handler returned unexpected body: got %v want %v", resWriter.Body.String(), tt.expectedText)

		})
	}
}
