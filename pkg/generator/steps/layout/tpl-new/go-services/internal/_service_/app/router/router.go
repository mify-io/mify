package router

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type routerModel struct {
	CoreInclude string // Include path to core package
}

func newRouterModel(ctx *gencontext.GenContext) routerModel {
	return routerModel{
		CoreInclude: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
	}
}
