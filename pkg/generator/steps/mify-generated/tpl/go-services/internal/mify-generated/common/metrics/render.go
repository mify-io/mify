package metrics

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	curPath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetCommonPath().Abs(),
		"metrics",
	)
	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(curPath, "client_metrics.go")),
		render.NewFile(ctx, path.Join(curPath, "metrics.go")),
		render.NewFile(ctx, path.Join(curPath, "request_metrics.go")),
	)
}
