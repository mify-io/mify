package apigateway

import (
	"context"

	"github.com/chebykinn/mify/pkg/generator/core"
)

type ApiGatewayStep struct {
}

func NewApiGatewayStep() ApiGatewayStep {
	return ApiGatewayStep{}
}

func (s ApiGatewayStep) Name() string {
	return "ApiGateway"
}

func (s ApiGatewayStep) ExecuteFunc() core.ExecuteFunc {
	return func(c *context.Context) *context.Context {
		return c
	}
}
