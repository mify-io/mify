package app

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type RequestExtraModel struct {
	CoreInclude string // Include path to core package
}

func newRequestExtraModel(ctx *gencontext.GenContext) RequestExtraModel {
	return RequestExtraModel{
		CoreInclude: ctx.GetWorkspace().GetCoreIncludePath(ctx.MustGetMifySchema().ServiceName),
	}
}
