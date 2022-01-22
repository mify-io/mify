package generator

import (
	"github.com/mify-io/mify/pkg/generator/core"
	apigateway "github.com/mify-io/mify/pkg/generator/steps/api-gateway"
	devrunner "github.com/mify-io/mify/pkg/generator/steps/dev-runner"
	layout "github.com/mify-io/mify/pkg/generator/steps/layout"
	openapi "github.com/mify-io/mify/pkg/generator/steps/openapi"
	schema "github.com/mify-io/mify/pkg/generator/steps/schema"
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
