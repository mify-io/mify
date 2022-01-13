package context

import (
	"fmt"
)

type SchemaContext struct {
	openapiSchemas OpenapiSchemas // service_name -> schema
}

func NewSchemaContext(openapiSchemas OpenapiSchemas) *SchemaContext {
	return &SchemaContext{
		openapiSchemas: openapiSchemas,
	}
}

func (c *SchemaContext) GetOpenapiSchemas(serviceName string) ServiceSchemas {
	schema, ok := c.openapiSchemas[serviceName]
	if !ok {
		panic(fmt.Sprintf("Schema for service '%s' wasn't found", serviceName))
	}

	return schema
}

func (c *SchemaContext) GetAllOpenapiSchemas() *OpenapiSchemas {
	return &c.openapiSchemas
}
