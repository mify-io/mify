package schema

import (
	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type SchemaStep struct {
}

func NewSchemaStep() SchemaStep {
	return SchemaStep{}
}

func (s SchemaStep) Name() string {
	return "schema-collector"
}

func (s SchemaStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	schemaCtx, err := execute(ctx)
	if err != nil {
		return core.Done, err
	}

	ctx.SetSchemaCtx(schemaCtx)
	return core.Done, nil
}
