package context

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenapiSchemas map[string]*openapi3.T

type SchemaContext struct {
	openapiSchemas OpenapiSchemas // service_name -> schema
}

func NewSchemaContext(openapiSchemas OpenapiSchemas) *SchemaContext {
	return &SchemaContext{
		openapiSchemas: openapiSchemas,
	}
}

func (c *SchemaContext) GetOpenapiSchemas() *OpenapiSchemas {
	return &c.openapiSchemas
}
