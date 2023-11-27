package app

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed mify_app.go.tpl
var mifyAppTemplate string

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	mifyAppModel := render.NewModel(ctx, newMifyAppModel(ctx))
	curPath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"app",
	)
	return render.RenderMany(templates,
		render.NewFile(ctx, path.Join(curPath, "mify_app.go")).SetModel(mifyAppModel),
	)
}
