package generator

import (
	"github.com/chebykinn/mify/pkg/generator/core"
	apigateway "github.com/chebykinn/mify/pkg/generator/steps/api-gateway"
	devrunner "github.com/chebykinn/mify/pkg/generator/steps/dev-runner"
	layout "github.com/chebykinn/mify/pkg/generator/steps/layout"
	openapi "github.com/chebykinn/mify/pkg/generator/steps/openapi"
	schema "github.com/chebykinn/mify/pkg/generator/steps/schema"
)

func BuildServicePipeline() core.Pipeline {
	return core.NewPipelineBuilder().
		Register(schema.NewSchemaStep()).
		Register(apigateway.NewApiGatewaySchemaStep()).
		Register(openapi.NewOpenapiStep()).
		Register(layout.NewLayoutStep()).
		Register(devrunner.NewDevRunnerStep()).
		Build()
}
