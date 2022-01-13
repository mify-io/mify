package apigateway

import (
	"github.com/chebykinn/mify/pkg/generator/core"
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
)

const (
	ApiGatewayName = "api-gateway" // TODO: it is not right place to declare it
)

type ApiGatewaySchemaStep struct {
}

func NewApiGatewaySchemaStep() ApiGatewaySchemaStep {
	return ApiGatewaySchemaStep{}
}

func (s ApiGatewaySchemaStep) Name() string {
	return "ApiGatewaySchema"
}

func (s ApiGatewaySchemaStep) Execute(ctx *gencontext.GenContext) (core.StepResult, error) {
	if ctx.GetServiceName() != ApiGatewayName {
		return core.Done, nil
	}

	foundPublicApis := scanPublicApis(ctx)
	res, err := updateApiGatewayOpenapiSchema(ctx, foundPublicApis)
	if err != nil {
		return core.Done, err
	}

	if res {
		return core.RepeatAll, nil
	}

	return core.Done, nil
}
