// TODO: split to smaller models

package tpl

import (
	"fmt"
	"path"

	"github.com/mify-io/mify/internal/mify/util"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

type GoServiceModel struct {
	Name string
}

type WorkspaceModel struct {
	Name       string
	BasePath   string
	GoRoot     string // Path to go_services
	Config     mifyconfig.WorkspaceConfig
	TplHeader  string
	GoServices []GoServiceModel
}

func NewWorkspaceModel(context *gencontext.GenContext) *WorkspaceModel {
	allSchemas := context.GetSchemaCtx().GetAllSchemas()
	goServices := make([]GoServiceModel, 0, len(*allSchemas))

	for serviceName, schemas := range *allSchemas {
		hasApi := schemas.GetOpenapi() != nil
		if hasApi && schemas.GetMify().Language == mifyconfig.ServiceLanguageGo {
			goServices = append(goServices, GoServiceModel{Name: serviceName})
		}
	}

	return &WorkspaceModel{
		Name:       context.GetWorkspace().Name,
		BasePath:   context.GetWorkspace().BasePath,
		GoRoot:     context.GetWorkspace().GetGoServicesPath(),
		Config:     context.GetWorkspace().Config,
		TplHeader:  "// THIS FILE IS AUTOGENERATED, DO NOT EDIT\n// Generated by mify",
		GoServices: goServices,
	}
}

// Path to include app.go
func (c WorkspaceModel) GetAppIncludePath(serviceName string) string {
	return fmt.Sprintf(
		"%s/go_services/internal/%s/generated/app",
		c.GetRepository(),
		serviceName)
}

func (c WorkspaceModel) GetApiSchemaDirAbsPath(serviceName string) string {
	return path.Join(c.BasePath, "schemas", serviceName, "api")
}

func (c WorkspaceModel) GetApiSchemaAbsPath(serviceName string) string {
	return path.Join(c.BasePath, "schemas", serviceName, "api/api.yaml")
}

// Path to api_generated.yaml
func (c WorkspaceModel) GetApiSchemaGenAbsPath(serviceName string) string {
	return path.Join(c.BasePath, "schemas", serviceName, "api/api_generated.yaml")
}

func (c *WorkspaceModel) GetRepository() string {
	return fmt.Sprintf("%s/%s/%s",
		c.Config.GitHost,
		c.Config.GitNamespace,
		c.Config.GitRepository)
}

// Name which can be used in generated go code
func (c WorkspaceModel) GetSafeName() string {
	return util.ToSafeGoVariableName(c.Name)
}
