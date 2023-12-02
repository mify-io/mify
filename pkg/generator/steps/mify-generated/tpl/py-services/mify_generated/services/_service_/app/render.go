package app

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	basePath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"app",
	)

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(basePath, "__init__.py")),
		render.NewFile(ctx, path.Join(basePath, "mify_app.py")),
		render.NewFile(ctx, path.Join(basePath, "server.py")),
	)
}
