package generated

import (
	_ "embed"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/libraries/generated/configs"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/libraries/generated/logs"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/libraries/generated/metrics"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed __init__.py.tpl
var initTemplate string

func Render(ctx *gencontext.GenContext) error {
	initModel := struct{}{}
	initPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedAbsPath(), "__init__.py")
	if err := render.RenderOrSkipTemplate(initTemplate, initModel, initPath); err != nil {
		return render.WrapError("__init__.py", err)
	}

	if err := configs.Render(ctx); err != nil {
		return err
	}
	if err := logs.Render(ctx); err != nil {
		return err
	}
	if err := metrics.Render(ctx); err != nil {
		return err
	}
	return nil
}
