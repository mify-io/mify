package gencontext

import (
	"context"
	"log"
	"os"

	api_gateway_context "github.com/chebykinn/mify/pkg/generator/steps/api-gateway/context"
	openapi_context "github.com/chebykinn/mify/pkg/generator/steps/openapi/context"
	schema_context "github.com/chebykinn/mify/pkg/generator/steps/schema/context"
	"github.com/chebykinn/mify/pkg/mifyconfig"
	"github.com/chebykinn/mify/pkg/workspace"
)

type GenContext struct {
	goContext context.Context
	Logger    *log.Logger

	serviceName   string
	workspace     workspace.Description
	serviceConfig mifyconfig.ServiceConfig

	// Step contexts
	schema     *schema_context.SchemaContext
	openapi    *openapi_context.OpenapiContext
	apiGateway *api_gateway_context.ApiGatewayContext
}

func NewGenContext(
	goContext context.Context,
	serviceName string,
	workspaceDescription workspace.Description,
	serviceConfig mifyconfig.ServiceConfig) *GenContext {

	return &GenContext{
		goContext:     goContext,
		Logger:        log.New(os.Stdout, "", 0),
		serviceName:   serviceName,
		workspace:     workspaceDescription,
		serviceConfig: serviceConfig,
	}
}

func (c *GenContext) GetGoContext() context.Context {
	return c.goContext
}

func (c *GenContext) GetServiceName() string {
	return c.serviceName
}

func (c *GenContext) GetWorkspace() *workspace.Description {
	return &c.workspace
}

func (c *GenContext) GetServiceConfig() *mifyconfig.ServiceConfig {
	return &c.serviceConfig
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
