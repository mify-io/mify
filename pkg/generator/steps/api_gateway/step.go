package apigateway

import (
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
)

type ApiGatewayStep struct {
}

func NewApiGatewayStep() ApiGatewayStep {
	return ApiGatewayStep{}
}

func (s ApiGatewayStep) Name() string {
	return "ApiGateway"
}

func (s ApiGatewayStep) Execute(ctx *gencontext.GenContext) error {
	return nil
}
