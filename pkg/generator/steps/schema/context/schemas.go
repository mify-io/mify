package context

import (
	"github.com/mify-io/mify/pkg/mifyconfig"
)

const (
	MainSchemaName      = "api.yaml"
	GeneratedSchemaName = "api_generated.yaml"
)

type AllSchemas map[string]*ServiceSchemas // service name -> schemas

type ServiceSchemas struct {
	openApi OpenapiServiceSchemas
	mify    *mifyconfig.ServiceConfig
}

func NewServiceSchemas(openApi OpenapiServiceSchemas, mify *mifyconfig.ServiceConfig) *ServiceSchemas {
	if mify == nil {
		panic("mify schema is required")
	}

	return &ServiceSchemas{
		openApi: openApi,
		mify:    mify,
	}
}

// Can be nill for some services
func (s ServiceSchemas) GetOpenapi() OpenapiServiceSchemas {
	return s.openApi
}

func (s ServiceSchemas) GetMify() *mifyconfig.ServiceConfig {
	return s.mify
}
