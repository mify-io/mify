package app

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed __init__.py.tpl
var initTemplate string

//go:embed request_extra.py.tpl
var requestExtraTemplate string

//go:embed service_extra.py.tpl
var serviceExtraTemplate string

func Render(ctx *gencontext.GenContext) error {
	initModel := struct{}{}
	initPath := ctx.GetWorkspace().GetPythonAppSubAbsPath(ctx.GetServiceName(), "__init__.py")
	if err := render.RenderOrSkipTemplate(initTemplate, initModel, initPath); err != nil {
		return render.WrapError("init", err)
	}

	reqExtraModel := newRequestExtraModel(ctx)
	reqExtraPath := ctx.GetWorkspace().GetPythonAppSubAbsPath(ctx.GetServiceName(), "request_extra.py")
	if err := render.RenderOrSkipTemplate(requestExtraTemplate, reqExtraModel, reqExtraPath); err != nil {
		return render.WrapError("request extra", err)
	}

	serviceExtraModel := newServiceExtraModel(ctx)
	serviceExtraPath := ctx.GetWorkspace().GetPythonAppSubAbsPath(ctx.GetServiceName(), "service_extra.py")
	if err := render.RenderOrSkipTemplate(serviceExtraTemplate, serviceExtraModel, serviceExtraPath); err != nil {
		return render.WrapError("service extra", err)
	}

	return nil
}
