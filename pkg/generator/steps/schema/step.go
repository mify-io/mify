package schema

import (
	"context"

	"github.com/chebykinn/mify/pkg/generator/core"
)

type SchemaStep struct {
}

func NewSchemaStep() SchemaStep {
	return SchemaStep{}
}

func (s SchemaStep) Name() string {
	return "Openapi"
}

func (s SchemaStep) ExecuteFunc() core.ExecuteFunc {
	return func(c *context.Context) *context.Context {
		return c
	}
}
