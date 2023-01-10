package generatedcore

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	dir := path.Join(ctx.GetWorkspace().GetJsServiceAbsPath(ctx.GetServiceName()), "generated", "core")

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(dir, "context.js")),
	)
}
