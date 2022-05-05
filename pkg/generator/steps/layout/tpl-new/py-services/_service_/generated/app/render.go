package app

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed __init__.py.tpl
var initTemplate string

//go:embed mify_app.py.tpl
var mifyAppTemplate string

//go:embed server.py.tpl
var serverTemplate string

func Render(ctx *gencontext.GenContext) error {
	initModel := struct{}{}
	initPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonGeneratedAppRelPath(ctx.GetServiceName()), "__init__.py")
	if err := render.RenderOrSkipTemplate(initTemplate, initModel, initPath); err != nil {
		return render.WrapError("__init__.py", err)
	}

	mifyAppModel := newMifyAppModel(ctx)
	mifyAppPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonGeneratedAppRelPath(ctx.GetServiceName()), "mify_app.py")
	if err := render.RenderTemplate(mifyAppTemplate, mifyAppModel, mifyAppPath); err != nil {
		return render.WrapError("mify_app.py", err)
	}

	serverModel, err := newServerModel(ctx)
	if err != nil {
		return err
	}
	serverPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonGeneratedAppRelPath(ctx.GetServiceName()), "server.py")
	if err := render.RenderTemplate(serverTemplate, serverModel, serverPath); err != nil {
		return render.WrapError("server.py", err)
	}

	return nil
}
