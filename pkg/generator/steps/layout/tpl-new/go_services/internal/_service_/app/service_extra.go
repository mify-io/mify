package app

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type ServiceExtraModel struct {
	CoreInclude string // Include path to core package
}

func newServiceExtraModel(ctx *gencontext.GenContext) ServiceExtraModel {
	return ServiceExtraModel{
		CoreInclude: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
	}
}
