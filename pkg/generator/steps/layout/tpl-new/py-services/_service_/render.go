package service

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/_service_/app"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed __init__.py.tpl
var initTemplate string

//go:embed __main__.py.tpl
var mainTemplate string

//go:embed Dockerfile.tpl
var dockerfileTemplate string

func Render(ctx *gencontext.GenContext) error {
	initModel := struct{}{}
	initPath := ctx.GetWorkspace().GetPythonServiceSubAbsPath(ctx.GetServiceName(), "__init__.py")
	if err := render.RenderOrSkipTemplate(initTemplate, initModel, initPath); err != nil {
		return render.WrapError("init", err)
	}

	mainModel := render.NewDefaultModel(ctx)
	mainPath := ctx.GetWorkspace().GetPythonServiceSubAbsPath(ctx.GetServiceName(), "__main__.py")
	if err := render.RenderOrSkipTemplate(mainTemplate, mainModel, mainPath); err != nil {
		return render.WrapError("main", err)
	}

	dockerfileModel := newDockerfileModel(ctx)
	dockerfilePath := ctx.GetWorkspace().GetPythonServiceSubAbsPath(ctx.GetServiceName(), "Dockerfile")
	if err := render.RenderOrSkipTemplate(dockerfileTemplate, dockerfileModel, dockerfilePath); err != nil {
		return render.WrapError("Dockerfile", err)
	}

	if err := app.Render(ctx); err != nil {
		return render.WrapError("app", err)
	}

	return nil
}
