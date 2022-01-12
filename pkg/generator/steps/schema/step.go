package schema

import gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"

type SchemaStep struct {
}

func NewSchemaStep() SchemaStep {
	return SchemaStep{}
}

func (s SchemaStep) Name() string {
	return "Openapi"
}

func (s SchemaStep) Execute(ctx *gencontext.GenContext) error {
	schemaCtx, err := execute(ctx)
	if err != nil {
		return err
	}

	ctx.SetSchemaCtx(schemaCtx)
	return nil
}
