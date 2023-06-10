package service

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/render"
	"github.com/mify-io/mify/pkg/workspace"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/service/tpl"
)

//go:embed tpl/api.yaml.tpl
var apiSchemaTemplate string

func CreateService(mutContext *mutators.MutatorContext, language mifyconfig.ServiceLanguage, template string, serviceName string) error {
	mutContext.GetLogger().Printf("Creating service '%s' ...", serviceName)

	err := validateLangAndTemplateForService(language, template)
	if err != nil {
		return err
	}

	return createServiceImpl(mutContext, language, serviceName, template, true)
}

func CreateFrontend(mutContext *mutators.MutatorContext, template string, name string) error {
	mutContext.GetLogger().Printf("Creating frontend '%s' ...", name)

	if template == "nuxtjs" || template == "react-ts" {
		return createServiceImpl(mutContext, mifyconfig.ServiceLanguageJs, name, template, false)
	}

	return fmt.Errorf("unknown template %s", template)
}

func CreateApiGateway(mutContext *mutators.MutatorContext) error {
	exists, err := checkServiceExists(mutContext, workspace.ApiGatewayName)
	if err != nil {
		return fmt.Errorf("can't check if service exists: %w", err)
	}

	if exists {
		return fmt.Errorf("api gateway already exists, skipping creation")
	}

	err = CreateService(mutContext, mifyconfig.ServiceLanguageGo, "", workspace.ApiGatewayName)
	if err != nil {
		return err
	}

	return nil
}

func createServiceImpl(
	mutContext *mutators.MutatorContext,
	language mifyconfig.ServiceLanguage,
	serviceName string,
	template string,
	addOpenApi bool) error {

	conf := mifyconfig.ServiceConfig{
		ServiceName: serviceName,
		Template:    template,
		Language:    language,
	}

	err := conf.Dump(mutContext.GetDescription().GetMifySchemaAbsPath(serviceName))
	if err != nil {
		return err
	}

	if addOpenApi {
		openapiSchemaPath := mutContext.GetDescription().GetApiSchemaAbsPath(serviceName, workspace.MainApiSchemaName)
		err := render.RenderTemplate(
			apiSchemaTemplate,
			tpl.NewApiSchemaModel(serviceName),
			openapiSchemaPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkServiceExists(mutContext *mutators.MutatorContext, serviceName string) (bool, error) {
	schemasDirAbsPath := mutContext.GetDescription().GetApiSchemaDirAbsPath(serviceName)
	if _, err := os.Stat(schemasDirAbsPath); os.IsNotExist(err) {
		return false, nil
	}

	files, err := os.ReadDir(schemasDirAbsPath)
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

func validateLangAndTemplateForService(language mifyconfig.ServiceLanguage, template string) error {
	switch language {
	case mifyconfig.ServiceLanguageJs:
		if template != "expressjs" {
			return fmt.Errorf("unsupported template '%s' for language '%s'", template, language)
		}
	}

	return nil
}
