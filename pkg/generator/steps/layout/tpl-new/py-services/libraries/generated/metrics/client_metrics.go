package metrics

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type clientMetricsModel struct {
	TplHeader string
}

func newClientMetricsModel(ctx *gencontext.GenContext) clientMetricsModel {
	return clientMetricsModel{
		TplHeader: ctx.GetWorkspace().TplHeaderPy,
	}
}
