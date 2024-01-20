package logs

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
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetCommonPath().Abs(),
		"logs",
	)
	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(basePath, "logger.py")),
	)
}
