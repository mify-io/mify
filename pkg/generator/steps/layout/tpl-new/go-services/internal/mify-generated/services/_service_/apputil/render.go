package apputil

import (
	"embed"
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed apputil.go.tpl
var appUtilTemplate string

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	appUtilModel := render.NewModel(ctx, newAppUtilModel(ctx))
	curPath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"apputil",
	)
	return render.RenderMany(templates,
		render.NewFile(ctx, path.Join(curPath, "apputil.go")).SetModel(appUtilModel),
	)
}
