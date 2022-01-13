package openapi

import (
	"github.com/chebykinn/mify/pkg/generator/core"
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
)

type OpenapiStep struct {
}

func NewOpenapiStep() OpenapiStep {
	return OpenapiStep{}
}

func (s OpenapiStep) Name() string {
	return "Openapi"
}

func (s OpenapiStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	return core.Done, nil
}
