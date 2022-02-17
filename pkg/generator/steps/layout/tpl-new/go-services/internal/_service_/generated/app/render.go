package app

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed mify_app.go.tpl
var mifyAppTemplate string

//go:embed server.go.tpl
var serverTemplate string

func Render(ctx *gencontext.GenContext) error {
	mifyAppModel := newMifyAppModel(ctx)
	mifyAppPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGeneratedAppRelPath(ctx.GetServiceName()), "mify_app.go")
	if err := render.RenderTemplate(mifyAppTemplate, mifyAppModel, mifyAppPath); err != nil {
		return render.WrapError("mify_app.go", err)
	}

	serverModel, err := newServerModel(ctx)
	if err != nil {
		return err
	}
	serverPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGeneratedAppRelPath(ctx.GetServiceName()), "server.go")
	if err := render.RenderTemplate(serverTemplate, serverModel, serverPath); err != nil {
		return render.WrapError("server.go", err)
	}

	return nil
}
