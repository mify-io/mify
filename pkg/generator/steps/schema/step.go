package schema

import "github.com/chebykinn/mify/pkg/generator/context"

type SchemaStep struct {
}

func NewSchemaStep() SchemaStep {
	return SchemaStep{}
}

func (s SchemaStep) Name() string {
	return "Openapi"
}

func (s SchemaStep) Execute(ctx *context.GenContext) error {
	schemaCtx, err := execute(ctx)
	if err != nil {
		return err
	}

	ctx.SetSchemaCtx(schemaCtx)
	return nil
}
