package service

import (
	_ "embed"
	"fmt"
	"io/ioutil"

	"github.com/chebykinn/mify/pkg/generator/templater"
	"github.com/chebykinn/mify/pkg/mifyconfig"
	"github.com/chebykinn/mify/pkg/workspace"
	"github.com/chebykinn/mify/pkg/workspace/mutators"
	"github.com/chebykinn/mify/pkg/workspace/mutators/service/tpl"
)

//go:embed tpl/api.yaml.tpl
var apiSchemaTemplate string

func CreateService(mutContext *mutators.MutatorContext, language mifyconfig.ServiceLanguage, serviceName string) error {
	fmt.Printf("Creating service: %s\n", serviceName)

	openapiSchemaPath := mutContext.GetDescription().GetApiSchemaAbsPath(serviceName, workspace.MainApiSchemaName)
	err := templater.RenderTemplate(
		"openapiSchema",
		apiSchemaTemplate,
		tpl.NewApiSchemaModel(serviceName),
		openapiSchemaPath)
	if err != nil {
		return err
	}

	conf := mifyconfig.ServiceConfig{
		ServiceName: serviceName,
		Language:    language,
	}

	err = conf.Dump(mutContext.GetDescription().GetMifySchemaAbsPath(serviceName))
	if err != nil {
		return err
	}

	return nil
}

func CreateFrontend(mutContext *mutators.MutatorContext, template string, name string) error {
	if template == "vue_js" {
		CreateService(mutContext, mifyconfig.ServiceLanguageJs, name)
		return nil
	}

	return fmt.Errorf("unknown template %s", template)
}

func TryCreateApiGateway(mutContext *mutators.MutatorContext) (bool, error) {
	exists, err := checkServiceExists(mutContext)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Printf("Api gateway already exists. Skipping creation... \n")
		return false, nil
	}

	err = CreateService(mutContext, mifyconfig.ServiceLanguageGo, workspace.ApiGatewayName)
	if err != nil {
		return false, err
	}

	return false, nil
}

func checkServiceExists(mutContext *mutators.MutatorContext) (bool, error) {
	schemasDirAbsPath := mutContext.GetDescription().GetApiSchemaDirAbsPath(workspace.ApiGatewayName)
	files, err := ioutil.ReadDir(schemasDirAbsPath)
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}
