package context

import (
	"fmt"

	"github.com/mify-io/mify/pkg/mifyconfig"
)

type SchemaContext struct {
	schemas AllSchemas
}

func NewSchemaContext(schemas AllSchemas) *SchemaContext {
	return &SchemaContext{
		schemas: schemas,
	}
}

func (c *SchemaContext) MustGetServiceSchemas(serviceName string) *ServiceSchemas {
	res := c.GetServiceSchemas(serviceName)
	if res == nil {
		panic(fmt.Sprintf("Schemas for service '%s' wasn't found", serviceName))
	}

	return res
}

func (c *SchemaContext) GetServiceSchemas(serviceName string) *ServiceSchemas {
	schemas, ok := c.schemas[serviceName]
	if !ok {
		return nil
	}

	return schemas
}

func (c *SchemaContext) GetAllSchemas() AllSchemas {
	return c.schemas
}

// Sugar

func (sc SchemaContext) GetMifySchema(serviceName string) *mifyconfig.ServiceConfig {
	schemas := sc.GetServiceSchemas(serviceName)
	if schemas == nil {
		return nil
	}

	return schemas.mify
}

func (sc SchemaContext) MustGetMifySchema(serviceName string) *mifyconfig.ServiceConfig {
	schemas := sc.MustGetServiceSchemas(serviceName)
	if schemas.mify == nil {
		panic(fmt.Sprintf("mify schema is missing for service '%s'", serviceName))
	}

	return schemas.mify
}
