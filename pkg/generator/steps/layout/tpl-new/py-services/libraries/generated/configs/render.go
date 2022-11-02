package configs

import (
	_ "embed"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed static.py.tpl
var staticTemplate string

//go:embed dynamic.py.tpl
var dynamicTemplate string

func Render(ctx *gencontext.GenContext) error {
	staticModel := newStaticModel(ctx)
	staticPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedConfigsAbsPath(), "static.py")
	if err := render.RenderOrSkipTemplate(staticTemplate, staticModel, staticPath); err != nil {
		return render.WrapError("static.py", err)
	}

	dynamicModel := newDynamicModel(ctx)
	dynamicPath := filepath.Join(ctx.GetWorkspace().GetPythonServicesLibrariesGeneratedConfigsAbsPath(), "dynamic.py")
	if err := render.RenderOrSkipTemplate(dynamicTemplate, dynamicModel, dynamicPath); err != nil {
		return render.WrapError("dynamic.py", err)
	}

	return nil
}
