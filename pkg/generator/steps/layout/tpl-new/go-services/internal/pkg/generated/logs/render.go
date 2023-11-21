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
	curPath := path.Join(
		ctx.GetWorkspace().GetGoServicesAbsPath(),
		"internal",
		"pkg",
		"generated",
		"logs",
	)

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(curPath, "logger.go")),
	)
}
