package generator

import (
	"github.com/mify-io/mify/pkg/generator/core"
	apigateway "github.com/mify-io/mify/pkg/generator/steps/api-gateway"
	devrunner "github.com/mify-io/mify/pkg/generator/steps/dev-runner"
	"github.com/mify-io/mify/pkg/generator/steps/migrate"
	"github.com/mify-io/mify/pkg/generator/steps/layout"
	mifygenerated "github.com/mify-io/mify/pkg/generator/steps/mify-generated"
	"github.com/mify-io/mify/pkg/generator/steps/openapi"
	"github.com/mify-io/mify/pkg/generator/steps/postgres"
	"github.com/mify-io/mify/pkg/generator/steps/prepare"
	"github.com/mify-io/mify/pkg/generator/steps/schema"
)

func BuildServicePipeline() core.Pipeline {
	return core.NewPipelineBuilder().
		Register(schema.NewSchemaStep()).
		Register(prepare.NewPrepareStep()).
		Register(apigateway.NewApiGatewaySchemaStep()).
		Register(migrate.NewMigrateStep()).
		Register(openapi.NewOpenapiStep()).
		Register(postgres.NewPostgresStep()).
		Register(layout.NewLayoutStep()).
		Register(mifygenerated.NewMifyGeneratedStep()).
		Register(devrunner.NewDevRunnerStep()).
		Build()
}
