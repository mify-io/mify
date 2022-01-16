package devrunner

import (
	_ "embed"
	"os"
	"path"

	"github.com/chebykinn/mify/internal/mify/util"
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"github.com/chebykinn/mify/pkg/generator/steps/dev-runner/tpl"
	"github.com/chebykinn/mify/pkg/generator/templater"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

const (
	DevRunnerName = "dev-runner"
)

//go:embed tpl/main.go.tpl
var devRunnerTemplate string

func execute(ctx *gencontext.GenContext) error {
	services := make([]tpl.TargetService, 0)
	for serviceName := range *ctx.GetSchemaCtx().GetAllOpenapiSchemas() {
		// TODO: temporary check if service is already generated. Remove in future
		path := ctx.GetWorkspace().GetAppPath(serviceName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		services = append(services, tpl.NewTargetService(
			serviceName,
			util.ToSafeGoVariableName(serviceName),
			ctx.GetWorkspace().GetAppIncludePath(serviceName),
		))
	}

	model := tpl.NewModel(ctx.GetWorkspace().TplHeader, services)
	err := templater.RenderTemplate("devRunner", devRunnerTemplate, model, buildPathToMainGo(ctx))
	if err != nil {
		return err
	}

	// TODO: move to service creation
	config := mifyconfig.ServiceConfig{
		ServiceName: DevRunnerName,
		Language:    mifyconfig.ServiceLanguageGo,
	}
	err = mifyconfig.SaveServiceConfig(ctx.GetWorkspace().BasePath, DevRunnerName, config)
	if err != nil {
		return err
	}

	return nil
}

func buildPathToMainGo(ctx *gencontext.GenContext) string {
	cmd := ctx.GetWorkspace().GetCmdPath(DevRunnerName)
	return path.Join(cmd, "main.go")
}
