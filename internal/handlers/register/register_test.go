package register

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"todo-go/internal/store"
	"todo-go/internal/store/mock"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewRegisterHandler(t *testing.T) {
	userStore := &mock.UserStoreMock{}
	type args struct {
		UserStore store.UserStore
	}
	tests := []struct {
		name string
		args args
		want *RegisterHandler
	}{
		{
			name: "default",
			args: args{UserStore: userStore},
			want: &RegisterHandler{userStore: userStore},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRegisterHandler(tt.args.UserStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRegisterHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterHandler_Get(t *testing.T) {

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       []byte
	}{
		{
			name:               "default",
			expectedStatusCode: http.StatusOK,
			expectedBody:       []byte("Create an account"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			handler := NewRegisterHandler(&mock.UserStoreMock{})
			req, _ := http.NewRequest("POST", "/", nil)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			handler.Get(rr, req)

			assert.Equal(tt.expectedStatusCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, tt.expectedStatusCode)

			assert.True(bytes.Contains(rr.Body.Bytes(), tt.expectedBody), "handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)

		})
	}
}

func TestRegisterHandler_Post(t *testing.T) {

	testCases := []struct {
		name               string
		email              string
		password           string
		createUserError    error
		expectedStatusCode int
		expectedBody       []byte
	}{
		{
			name:               "success",
			email:              "test@example.com",
			password:           "password",
			expectedStatusCode: http.StatusOK,
			expectedBody:       []byte(`Registration successful`),
		},
		{
			name:               "fail - error creating user",
			email:              "test@example.com",
			password:           "password",
			createUserError:    gorm.ErrDuplicatedKey,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       []byte(`There was an error registering your account`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			userStore := &mock.UserStoreMock{}

			userStore.On("CreateUser", tc.email, tc.password).Return(tc.createUserError)

			handler := NewRegisterHandler(userStore)

			body := bytes.NewBufferString("email=" + tc.email + "&password=" + tc.password)

			req, _ := http.NewRequest("POST", "/", body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()

			handler.Post(rr, req)

			assert.Equal(tc.expectedStatusCode, rr.Code, "handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatusCode)

			assert.True(bytes.Contains(rr.Body.Bytes(), tc.expectedBody), "handler returned unexpected body: got %v want %v", rr.Body.String(), tc.expectedBody)

			userStore.AssertExpectations(t)
		})
	}
}
