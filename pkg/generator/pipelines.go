package generator

import (
	"github.com/chebykinn/mify/pkg/generator/core"
	apigateway "github.com/chebykinn/mify/pkg/generator/steps/api_gateway"
	openapi "github.com/chebykinn/mify/pkg/generator/steps/openapi"
	schema "github.com/chebykinn/mify/pkg/generator/steps/schema"
)

func BuildServicePipeline() core.Pipeline {
	return core.NewPipelineBuilder().
		Register(schema.NewSchemaStep()).
		Register(openapi.NewOpenapiStep()).
		Register(apigateway.NewApiGatewayStep()).
		Build()
}
