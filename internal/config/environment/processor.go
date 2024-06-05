package environment

import "github.com/kelseyhightower/envconfig"

type EnvConfigProcessor interface {
	Process(prefix string, spec interface{}) error
}

type DefaultEnvConfigProcessor struct{}

func (d DefaultEnvConfigProcessor) Process(prefix string, spec interface{}) error {
	return envconfig.Process(prefix, spec)
}
