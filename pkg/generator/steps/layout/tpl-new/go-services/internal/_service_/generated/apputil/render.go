package apputil

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed apputil.go.tpl
var appUtilTemplate string

func Render(ctx *gencontext.GenContext) error {
	appUtilModel := newAppUtilModel(ctx)
	appUtilPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGeneratedRelPath(ctx.GetServiceName()), "apputil/apputil.go")
	if err := render.RenderTemplate(appUtilTemplate, appUtilModel, appUtilPath); err != nil {
		return render.WrapError("apputil/apputil.go", err)
	}

	return nil
}
