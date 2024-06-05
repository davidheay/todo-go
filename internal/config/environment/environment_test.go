package environment

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"todo-go/internal/config"
	mok "todo-go/internal/config/environment/mock"
)

func TestNewEnvConfig(t *testing.T) {
	processor := &mok.ProcessorMock{}
	envConfig := NewEnvConfig(processor)

	if envConfig.EnvConfigProcessor != processor {
		t.Errorf("EnvConfigProcessor is not equal, want %v, got %v", processor, envConfig.EnvConfigProcessor)
	}
}

func TestMustLoadConfig(t *testing.T) {
	os.Clearenv()
	defer os.Clearenv()

	tests := []struct {
		name  string
		env   map[string]string
		mock  EnvConfigProcessor
		want  *config.Config
		error error
	}{
		{
			name: "default",
			env:  map[string]string{},
			mock: DefaultEnvConfigProcessor{},
			want: &config.Config{
				Port:              ":8080",
				DatabaseName:      "todogo",
				SessionCookieName: "session",
			},
		},
		{
			name: "custom port",
			env: map[string]string{
				"PORT": ":3000",
			},
			mock: DefaultEnvConfigProcessor{},
			want: &config.Config{
				Port:              ":3000",
				DatabaseName:      "todogo",
				SessionCookieName: "session",
			},
		},
		{
			name: "custom database name",
			env: map[string]string{
				"DATABASE_NAME": "test_db",
			},
			mock: DefaultEnvConfigProcessor{},
			want: &config.Config{
				Port:              ":8080",
				DatabaseName:      "test_db",
				SessionCookieName: "session",
			},
		},
		{
			name: "custom session cookie name",
			env: map[string]string{
				"SESSION_COOKIE_NAME": "test_session",
			},
			mock: DefaultEnvConfigProcessor{},
			want: &config.Config{
				Port:              ":8080",
				DatabaseName:      "todogo",
				SessionCookieName: "test_session",
			},
		},
		{
			name: "error env",
			mock: &mok.ProcessorMock{},
			env:  map[string]string{},
			want: &config.Config{
				Port:              ":8080",
				DatabaseName:      "todogo",
				SessionCookieName: "test_session",
			},
			error: fmt.Errorf("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				if err := os.Setenv(k, v); err != nil {
					t.Fatal(err)
				}
			}
			envConfig := EnvConfig{EnvConfigProcessor: tt.mock}
			if tt.error == nil {
				cfg := envConfig.MustLoadConfig()
				if !reflect.DeepEqual(cfg, tt.want) {
					t.Errorf("MustLoadConfig() = %v, want %v", cfg, tt.want)
				}
			} else {
				defer func() {
					r := recover()
					if r == nil {
						t.Errorf("MustLoadConfig() should panic")
					}
				}()
				envConfig.MustLoadConfig()
			}
			os.Clearenv()
		})
	}
}
