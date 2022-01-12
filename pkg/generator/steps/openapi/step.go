package openapi

import (
	generator "github.com/chebykinn/mify/pkg/generator/context"
)

type OpenapiStep struct {
}

func NewOpenapiStep() OpenapiStep {
	return OpenapiStep{}
}

func (s OpenapiStep) Name() string {
	return "Openapi"
}

func (s OpenapiStep) Execute(ctx *generator.GenContext) error {
	return nil
}
