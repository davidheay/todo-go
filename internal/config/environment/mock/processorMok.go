package mok

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

type ProcessorMock struct {
	mock.Mock
}

func (p *ProcessorMock) Process(prefix string, spec interface{}) error {
	return fmt.Errorf("error")
}
