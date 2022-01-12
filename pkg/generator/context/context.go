package context

import (
	api_gateway_context "github.com/chebykinn/mify/pkg/generator/steps/api_gateway/context"
	openapi_context "github.com/chebykinn/mify/pkg/generator/steps/openapi/context"
	schema_context "github.com/chebykinn/mify/pkg/generator/steps/schema/context"
)

type Context struct {
	schema     *schema_context.SchemaContext
	openapi    *openapi_context.OpenapiContext
	apiGateway *api_gateway_context.ApiGatewayContext
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) GetSchemaCtx() *schema_context.SchemaContext {
	if c.schema == nil {
		panic("Schema context is not filled")
	}
	return c.schema
}

func (c *Context) GetOpenapiCtx() *openapi_context.OpenapiContext {
	if c.openapi == nil {
		panic("Openapi context is not filled")
	}
	return c.openapi
}

func (c *Context) GetApiGatewayCtx() *api_gateway_context.ApiGatewayContext {
	if c.apiGateway == nil {
		panic("Api gateway context is not filled")
	}
	return c.apiGateway
}
