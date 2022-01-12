package openapi

import (
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

func (s OpenapiStep) Execute(ctx *gencontext.GenContext) error {
	return nil
}
