package openapi

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed app.py.tpl
var appTemplate string

func Render(ctx *gencontext.GenContext) error {
	appModel := newAppModel(ctx)
	appPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonServiceGeneratedOpenAPIRelPath(ctx.GetServiceName()), "app.py")
	if err := render.RenderTemplate(appTemplate, appModel, appPath); err != nil {
		return render.WrapError("app.py", err)
	}

	return nil
}
