package config

type Config struct {
	Port              string `envconfig:"PORT" default:":8080"`
	DatabaseName      string `envconfig:"DATABASE_NAME" default:"todogo"`
	DatabasePassword  string `envconfig:"DATABASE_PASSWORD"`
	SessionCookieName string `envconfig:"SESSION_COOKIE_NAME" default:"session"`
}
type ConfigProcessor interface {
	MustLoadConfig() *Config
}
