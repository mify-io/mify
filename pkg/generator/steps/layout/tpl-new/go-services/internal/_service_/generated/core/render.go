package core

import (
	_ "embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed request_context.go.tpl
var requestContextTemplate string

//go:embed service_context.go.tpl
var serviceContextTemplate string

//go:embed helpers.go.tpl
var helpersTemplate string

func Render(ctx *gencontext.GenContext) error {
	requestContextModel := newRequestContextModel(ctx)
	requestContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "request_context.go")
	if err := render.RenderTemplate(requestContextTemplate, requestContextModel, requestContextPath); err != nil {
		return render.WrapError("request_context", err)
	}

	serviceContextModel := newServiceContextModel(ctx)
	serviceContextPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "service_context.go")
	if err := render.RenderTemplate(serviceContextTemplate, serviceContextModel, serviceContextPath); err != nil {
		return render.WrapError("service_context", err)
	}

	helpersModel := newHelpersModel(ctx)
	helpersPath := path.Join(ctx.GetWorkspace().BasePath, ctx.GetWorkspace().GetGoServiceGeneratedCoreRelPath(ctx.GetServiceName()), "helpers.go")
	if err := render.RenderTemplate(helpersTemplate, helpersModel, helpersPath); err != nil {
		return render.WrapError("helpers", err)
	}

	return nil
}
