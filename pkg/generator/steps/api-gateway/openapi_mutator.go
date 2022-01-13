package apigateway

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/generator/steps/schema/context"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
)

func updateApiGatewayOpenapiSchema(genContex *gencontext.GenContext, publicApis PublicApis) (bool, error) {
	currentGeneratedSchema := genContex.GetSchemaCtx().GetOpenapiSchemas(genContex.GetServiceName()).GetGeneratedSchema()

	if len(publicApis) == 0 {
		if currentGeneratedSchema == nil {
			return false, nil
		}

		return true, removeGeneratedSchema(genContex)
	}

	newSchema := buildOpenapiSchema(publicApis)
	schemasEqual, err := areSchemasEqual(currentGeneratedSchema, newSchema)
	if err != nil {
		return false, err
	}

	if schemasEqual {
		return false, nil
	}

	newSchemaYaml, err := marshalYaml(newSchema)
	if err != nil {
		return false, err
	}

	err = ioutil.WriteFile(getAbsPathToGeneratedApiSchema(genContex), []byte(newSchemaYaml), 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getAbsPathToGeneratedApiSchema(genContex *gencontext.GenContext) string {
	serviceName := genContex.GetServiceName()
	return genContex.GetWorkspace().GetApiSchemaAbsPath(serviceName, context.GeneratedSchemaName)
}

func removeGeneratedSchema(genContex *gencontext.GenContext) error {
	return os.Remove(getAbsPathToGeneratedApiSchema(genContex))
}

func buildOpenapiSchema(publicApis PublicApis) *openapi3.T {
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Version:     "1.0.0",
			Title:       "api-gateway",
			Description: "Generated by mify",
		},
		Paths: make(openapi3.Paths),
	}

	for _, apis := range publicApis {
		for path, api := range apis {
			doc.Paths[path] = api
		}
	}

	return doc
}

func areSchemasEqual(schema1 *openapi3.T, schema2 *openapi3.T) (bool, error) {
	if (schema1 == nil || schema2 == nil) && (schema1 != schema2) {
		return false, nil
	}

	res1, err := marshalYaml(schema1)
	if err != nil {
		return false, err
	}

	res2, err := marshalYaml(schema2)
	if err != nil {
		return false, err
	}

	return res1 == res2, nil
}

// Marshal yaml is not implemented by kin-openapi. So a little hack
func marshalYaml(schema *openapi3.T) (string, error) {
	jsonBuf, err := schema.MarshalJSON()
	if err != nil {
		return "", err
	}

	tmp := map[string]interface{}{}
	err = json.Unmarshal(jsonBuf, &tmp)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(tmp); err != nil {
		return "", err
	}

	return buf.String(), nil
}
