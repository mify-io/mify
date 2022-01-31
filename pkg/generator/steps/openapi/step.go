package openapi

import (
	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type OpenapiStep struct {
}

func NewOpenapiStep() OpenapiStep {
	return OpenapiStep{}
}

func (s OpenapiStep) Name() string {
	return "openapi"
}

func (s OpenapiStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if ctx.GetMifySchema() == nil {
		return core.Done, nil // Some services (like dev-runner) could not have any scheme
	}

	if err := generateServiceOpenAPI(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}
