package context

import "github.com/getkin/kin-openapi/openapi3"

type OpenapiServiceSchemas map[string]*openapi3.T // schema name -> openapi schema

func (s OpenapiServiceSchemas) GetMainSchema() *openapi3.T {
	res, ok := s[MainSchemaName]
	if !ok {
		return nil
	}

	return res
}

func (s OpenapiServiceSchemas) GetGeneratedSchema() *openapi3.T {
	res, ok := s[GeneratedSchemaName]
	if !ok {
		return nil
	}

	return res
}
