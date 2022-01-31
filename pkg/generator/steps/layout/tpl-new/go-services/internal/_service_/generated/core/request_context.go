package core

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type requestContextModel struct {
	TplHeader         string
	MetricsImportPath string
}

func newRequestContextModel(ctx *gencontext.GenContext) requestContextModel {
	// TODO: move paths to description
	return requestContextModel{
		TplHeader:         ctx.GetWorkspace().TplHeader,
		MetricsImportPath: fmt.Sprintf("%s/internal/pkg/generated/metrics", ctx.GetWorkspace().GetGoModule()),
	}
}
