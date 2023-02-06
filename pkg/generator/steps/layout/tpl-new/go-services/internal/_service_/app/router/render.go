package router

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed router.go.tpl
var routerTemplate string

func Render(ctx *gencontext.GenContext) error {
	routerModel := newRouterModel(ctx)
	routerPath := ctx.GetWorkspace().GetAppSubAbsPath(ctx.GetServiceName(), "router/router.go")
	if err := render.RenderOrSkipTemplate(routerTemplate, routerModel, routerPath); err != nil {
		return render.WrapError("request extra", err)
	}

	return nil
}
