package gencontext

import (
	"context"

	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
	api_gateway_context "github.com/mify-io/mify/pkg/generator/steps/api-gateway/context"
	openapi_context "github.com/mify-io/mify/pkg/generator/steps/openapi/context"
	schema_context "github.com/mify-io/mify/pkg/generator/steps/schema/context"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/workspace"
	"go.uber.org/zap"
)

type GenContext struct {
	goContext          context.Context
	Logger             *zap.SugaredLogger
	executePoolFactory *ExecutePoolFactory

	serviceName string
	workspace   workspace.Description
	migrate     bool

	// Step contexts
	schema     *schema_context.SchemaContext
	openapi    *openapi_context.OpenapiContext
	apiGateway *api_gateway_context.ApiGatewayContext

	// libs
	EndpointsResolver *endpoints.EndpointsResolver
	vcsIntegration    *VcsIntegration
}

func NewGenContext(
	goContext context.Context,
	serviceName string,
	workspaceDescription workspace.Description,
	migrate bool) (*GenContext, error) {

	logger := initLogger(workspaceDescription.GetLogsDirectory())
	vcs, err := initVcsIntegration(workspaceDescription.BasePath)
	if err != nil {
		return nil, err
	}

	return &GenContext{
		goContext:         goContext,
		Logger:            logger.Sugar(),
		serviceName:       serviceName,
		workspace:         workspaceDescription,
		migrate:           migrate,
		EndpointsResolver: endpoints.NewEndpointsResolver(&workspaceDescription),
		vcsIntegration:    vcs,
	}, nil
}

func (c *GenContext) GetGoContext() context.Context {
	return c.goContext
}

func (c *GenContext) GetExecutePoolFactory() *ExecutePoolFactory {
	return c.executePoolFactory
}

func (c *GenContext) GetServiceName() string {
	return c.serviceName
}

func (c *GenContext) GetWorkspace() *workspace.Description {
	return &c.workspace
}

func (c *GenContext) GetMigrate() bool {
	return c.migrate
}

func (c *GenContext) GetVcsIntegration() *VcsIntegration {
	return c.vcsIntegration
}

func (c *GenContext) GetSchemaCtx() *schema_context.SchemaContext {
	if c.schema == nil {
		panic("Schema context is not filled")
	}
	return c.schema
}

func (c *GenContext) SetSchemaCtx(ctx *schema_context.SchemaContext) {
	if c.schema != nil {
		panic("Schema context is already filled")
	}
	c.schema = ctx
}

func (c *GenContext) GetOpenapiCtx() *openapi_context.OpenapiContext {
	if c.openapi == nil {
		panic("Openapi context is not filled")
	}
	return c.openapi
}

func (c *GenContext) SetOpenapiCtx(ctx *openapi_context.OpenapiContext) {
	if c.openapi != nil {
		panic("Openapi context is already filled")
	}
	c.openapi = ctx
}

func (c *GenContext) GetApiGatewayCtx() *api_gateway_context.ApiGatewayContext {
	if c.apiGateway == nil {
		panic("Api gateway context is not filled")
	}
	return c.apiGateway
}

func (c *GenContext) SetApiGatewayCtx(ctx *api_gateway_context.ApiGatewayContext) {
	if c.apiGateway == nil {
		panic("Api gateway context is not filled")
	}
	c.apiGateway = ctx
}

// Sugar

func (c *GenContext) GetMifySchema() *mifyconfig.ServiceConfig {
	return c.GetSchemaCtx().GetMifySchema(c.GetServiceName())
}

func (c *GenContext) MustGetMifySchema() *mifyconfig.ServiceConfig {
	return c.GetSchemaCtx().MustGetMifySchema(c.GetServiceName())
}
