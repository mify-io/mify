package devrunner

import (
	_ "embed"
	"os"
	"path"

	"github.com/mify-io/mify/internal/mify/util"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/dev-runner/tpl"
	"github.com/mify-io/mify/pkg/mifyconfig"
	"github.com/mify-io/mify/pkg/util/render"
	"github.com/mify-io/mify/pkg/workspace"
)

//go:embed tpl/main.go.tpl
var devRunnerTemplate string

func execute(ctx *gencontext.GenContext) error {
	services := make([]tpl.TargetService, 0)
	for serviceName, schemas := range ctx.GetSchemaCtx().GetAllSchemas() {
		if schemas.GetMify().Language != mifyconfig.ServiceLanguageGo {
			continue
		}

		// TODO: temporary check if service is already generated. Force regenerate in future
		path := ctx.GetWorkspace().GetGeneratedAppPath(serviceName)
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
	err := render.RenderTemplate(devRunnerTemplate, model, buildPathToMainGo(ctx))
	if err != nil {
		return err
	}

	return nil
}

func buildPathToMainGo(ctx *gencontext.GenContext) string {
	cmd := ctx.GetWorkspace().GetCmdAbsPath(workspace.DevRunnerName)
	return path.Join(cmd, "main.go")
}
