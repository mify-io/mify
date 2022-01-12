package apigateway

import (
	generator "github.com/chebykinn/mify/pkg/generator/context"
)

type ApiGatewayStep struct {
}

func NewApiGatewayStep() ApiGatewayStep {
	return ApiGatewayStep{}
}

func (s ApiGatewayStep) Name() string {
	return "ApiGateway"
}

func (s ApiGatewayStep) Execute(ctx *generator.GenContext) error {
	return nil
}
