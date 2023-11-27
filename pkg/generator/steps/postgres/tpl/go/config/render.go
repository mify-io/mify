package config

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	postgresConfigModel := NewPostgresConfigModel(ctx)
	basePath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetServicePath().Abs(),
		"postgres",
	)

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(basePath, "config.go")).SetModel(render.NewModel(ctx, postgresConfigModel)),
	)
}
