package service

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/mify-io/mify/pkg/generator/lib/endpoints"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/render"
	"github.com/mify-io/mify/pkg/workspace"
	"github.com/mify-io/mify/pkg/workspace/mutators"
	"github.com/mify-io/mify/pkg/workspace/mutators/service/tpl"
	"github.com/samber/lo"
)

//go:embed tpl/api.yaml.tpl
var apiSchemaTemplate string

func CreateService(mutContext *mutators.MutatorContext, language mifyconfig.ServiceLanguage, template string, serviceName string) (mifyconfig.ServiceConfig, error) {
	mutContext.GetLogger().Printf("Creating service '%s' ...", serviceName)

	err := validateLangAndTemplateForService(language, template)
	if err != nil {
		return mifyconfig.ServiceConfig{}, err
	}

	return createServiceImpl(mutContext, language, serviceName, template, true)
}

func CreateFrontend(mutContext *mutators.MutatorContext, template string, name string) error {
	mutContext.GetLogger().Printf("Creating frontend '%s' ...", name)

	templates := []string{
		"nuxtjs",
		"react-ts",
		"react-ts-nginx",
	}
	if lo.Contains(templates, template) {
		_, err := createServiceImpl(mutContext, mifyconfig.ServiceLanguageJs, name, template, false)
		return err
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

	_, err = CreateService(mutContext, mifyconfig.ServiceLanguageGo, "", workspace.ApiGatewayName)
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
	addOpenApi bool) (mifyconfig.ServiceConfig, error) {
	if language == mifyconfig.ServiceLanguagePython {
		serviceName = endpoints.SanitizeServiceName(serviceName)
	}

	conf := mifyconfig.ServiceConfig{
		ServiceName: serviceName,
		Template:    template,
		Language:    language,
	}

	err := conf.Dump(mutContext.GetDescription().GetMifySchemaAbsPath(serviceName))
	if err != nil {
		return mifyconfig.ServiceConfig{}, err
	}

	if addOpenApi {
		openapiSchemaPath := mutContext.GetDescription().GetApiSchemaAbsPath(serviceName, workspace.MainApiSchemaName)
		err := render.RenderTemplate(
			apiSchemaTemplate,
			tpl.NewApiSchemaModel(serviceName),
			openapiSchemaPath)
		if err != nil {
			return mifyconfig.ServiceConfig{}, err
		}
	}

	return conf, nil
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
