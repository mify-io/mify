package generator

import (
	"github.com/mify-io/mify/pkg/generator/core"
	apigateway "github.com/mify-io/mify/pkg/generator/steps/api-gateway"
	devrunner "github.com/mify-io/mify/pkg/generator/steps/dev-runner"
	"github.com/mify-io/mify/pkg/generator/steps/layout"
	"github.com/mify-io/mify/pkg/generator/steps/openapi"
	"github.com/mify-io/mify/pkg/generator/steps/postgres"
	"github.com/mify-io/mify/pkg/generator/steps/schema"
)

func BuildServicePipeline() core.Pipeline {
	return core.NewPipelineBuilder().
		Register(schema.NewSchemaStep()).
		Register(apigateway.NewApiGatewaySchemaStep()).
		Register(openapi.NewOpenapiStep()).
		Register(postgres.NewPostgresStep()).
		Register(layout.NewLayoutStep()).
		Register(devrunner.NewDevRunnerStep()).
		Build()
}
