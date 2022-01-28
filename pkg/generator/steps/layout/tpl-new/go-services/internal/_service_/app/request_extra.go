package app

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type requestExtraModel struct {
	CoreInclude string // Include path to core package
}

func newRequestExtraModel(ctx *gencontext.GenContext) requestExtraModel {
	return requestExtraModel{
		CoreInclude: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
	}
}
