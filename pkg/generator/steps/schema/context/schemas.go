package context

import "github.com/getkin/kin-openapi/openapi3"

const (
	MainSchemaName      = "api.yaml"
	GeneratedSchemaName = "api_generated.yaml"
)

type ServiceSchemas map[string]*openapi3.T    // schema name -> schema
type OpenapiSchemas map[string]ServiceSchemas // service name -> schemas

func (s ServiceSchemas) GetMainSchema() *openapi3.T {
	res, ok := s[MainSchemaName]
	if !ok {
		return nil
	}

	return res
}

func (s ServiceSchemas) GetGeneratedSchema() *openapi3.T {
	res, ok := s[GeneratedSchemaName]
	if !ok {
		return nil
	}

	return res
}
