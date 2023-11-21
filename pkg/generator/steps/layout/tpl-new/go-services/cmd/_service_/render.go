package service

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
		"cmd",
		ctx.GetServiceName(),
	)

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(curPath, "main.go")),
		render.NewFile(ctx, path.Join(curPath, "Dockerfile")),
	)
}
