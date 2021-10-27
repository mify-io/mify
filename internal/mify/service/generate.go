package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/service/generate"
)

const (
	apiSchemaPath = "schemas/%s/api"
	apiServicePath = "backend/internal/%s"
	svcLanguage generate.GeneratorLanguage = generate.GENERATOR_LANGUAGE_GO
)

func Generate(name string) error {
	_, err := config.ReadWorkspaceConfig()
	if err != nil {
		return err
	}
	// check if service exists

	if err := generateServiceOpenAPI(name); err != nil {
		return err
	}
	return nil
}

func generateServiceOpenAPI(name string) error {
	schemaPath := fmt.Sprintf(apiSchemaPath, name)
	if _, err := os.Stat(schemaPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("debug: skipping openapi generating, schema not found for service: %s\n", name)
		return nil
	}
	openapigen := generate.NewOpenAPIGenerator(schemaPath, svcLanguage)
	return openapigen.GenerateServer(fmt.Sprintf(apiServicePath, name))
}
