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
	serverModel := render.NewModel(ctx, newServerModel(ctx))
	curPath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"app",
	)
	return render.RenderMany(templates,
		render.NewFile(ctx, path.Join(curPath, "server.go")).SetModel(serverModel),
	)
}
