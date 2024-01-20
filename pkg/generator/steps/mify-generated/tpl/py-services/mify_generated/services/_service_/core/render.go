package core

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	clientsModel := newClientsModel(ctx)
	basePath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"core",
	)

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(basePath, "__init__.py")),
		render.NewFile(ctx, path.Join(basePath, "request_context.py")),
		render.NewFile(ctx, path.Join(basePath, "service_context.py")),
		render.NewFile(ctx, path.Join(basePath, "clients.py")).
			SetModel(render.NewModel(ctx, clientsModel)),
	)
}
