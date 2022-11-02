package metrics

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type requestMetricsModel struct {
	TplHeader      string
}

func newRequestMetricsModel(ctx *gencontext.GenContext) requestMetricsModel {
	return requestMetricsModel{
		TplHeader:      ctx.GetWorkspace().TplHeaderPy,
	}
}
