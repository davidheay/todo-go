package environment

import (
	"todo-go/internal/config"
)

type EnvConfig struct {
	EnvConfigProcessor
}

func NewEnvConfig(envConfigProcessor EnvConfigProcessor) *EnvConfig {
	return &EnvConfig{envConfigProcessor}
}

func (e *EnvConfig) MustLoadConfig() *config.Config {
	cfg, err := loadConfig(e.EnvConfigProcessor)
	if err != nil {
		panic(err)
	}
	return cfg
}

func loadConfig(processor EnvConfigProcessor) (*config.Config, error) {
	cfg := &config.Config{}
	if err := processor.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
