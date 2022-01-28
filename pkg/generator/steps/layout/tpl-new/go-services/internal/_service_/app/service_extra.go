package app

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type serviceExtraModel struct {
	CoreInclude string // Include path to core package
}

func newServiceExtraModel(ctx *gencontext.GenContext) serviceExtraModel {
	return serviceExtraModel{
		CoreInclude: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
	}
}
