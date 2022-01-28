package app

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed request_extra.go.tpl
var requestExtraTemplate string

//go:embed service_extra.go.tpl
var serviceExtraTemplate string

func Render(ctx *gencontext.GenContext) error {
	reqExtraModel := newRequestExtraModel(ctx)
	reqExtraPath := ctx.GetWorkspace().GetAppSubAbsPath(ctx.GetServiceName(), "request_extra.go")
	if err := render.RenderOrSkipTemplate(requestExtraTemplate, reqExtraModel, reqExtraPath); err != nil {
		return render.WrapError("request extra", err)
	}

	serviceExtraModel := newServiceExtraModel(ctx)
	serviceExtraPath := ctx.GetWorkspace().GetAppSubAbsPath(ctx.GetServiceName(), "service_extra.go")
	if err := render.RenderOrSkipTemplate(serviceExtraTemplate, serviceExtraModel, serviceExtraPath); err != nil {
		return render.WrapError("service extra", err)
	}

	return nil
}
