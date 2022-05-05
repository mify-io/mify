package logs

import (
	_ "embed"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed logger.py.tpl
var loggerTemplate string

func Render(ctx *gencontext.GenContext) error {
	loggerModel := newLoggerModel(ctx)
	loggerPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedLogsAbsPath(), "logger.py")
	if err := render.RenderOrSkipTemplate(loggerTemplate, loggerModel, loggerPath); err != nil {
		return render.WrapError("logger.py", err)
	}
	return nil
}
