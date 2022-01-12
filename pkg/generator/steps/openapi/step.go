package openapi

import (
	"context"

	"github.com/chebykinn/mify/pkg/generator/core"
)

type OpenapiStep struct {
}

func NewOpenapiStep() OpenapiStep {
	return OpenapiStep{}
}

func (s OpenapiStep) Name() string {
	return "Openapi"
}

func (s OpenapiStep) ExecuteFunc() core.ExecuteFunc {
	return func(c *context.Context) *context.Context {
		return c
	}
}
