package schema

import (
	"github.com/chebykinn/mify/pkg/generator/core"
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
)

type SchemaStep struct {
}

func NewSchemaStep() SchemaStep {
	return SchemaStep{}
}

func (s SchemaStep) Name() string {
	return "Openapi"
}

func (s SchemaStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	schemaCtx, err := execute(ctx)
	if err != nil {
		return core.Done, err
	}

	ctx.SetSchemaCtx(schemaCtx)
	return core.Done, nil
}
