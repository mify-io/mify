package core

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed __init__.py.tpl
var initTemplate string

//go:embed request_context.py.tpl
var requestContextTemplate string

//go:embed service_context.py.tpl
var serviceContextTemplate string

//go:embed clients.py.tpl
var clientsTemplate string

func Render(ctx *gencontext.GenContext) error {
	initModel := struct{}{}
	initPath := ctx.GetWorkspace().GetPythonAppSubAbsPath(ctx.GetServiceName(), "__init__.py")
	if err := render.RenderOrSkipTemplate(initTemplate, initModel, initPath); err != nil {
		return render.WrapError("__init__.py", err)
	}

	requestContextModel := newRequestContextModel(ctx)
	requestContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonServiceGeneratedCoreRelPath(ctx.GetServiceName()), "request_context.py")
	if err := render.RenderTemplate(requestContextTemplate, requestContextModel, requestContextPath); err != nil {
		return render.WrapError("request_context.py", err)
	}

	serviceContextModel := newServiceContextModel(ctx)
	serviceContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonServiceGeneratedCoreRelPath(ctx.GetServiceName()), "service_context.py")
	if err := render.RenderTemplate(serviceContextTemplate, serviceContextModel, serviceContextPath); err != nil {
		return render.WrapError("service_context.py", err)
	}

	clientsModel := newClientsModel(ctx)
	clientsPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetPythonServiceGeneratedCoreRelPath(ctx.GetServiceName()), "clients.py")
	if err := render.RenderTemplate(clientsTemplate, clientsModel, clientsPath); err != nil {
		return render.WrapError("clients.py", err)
	}

	return nil
}
