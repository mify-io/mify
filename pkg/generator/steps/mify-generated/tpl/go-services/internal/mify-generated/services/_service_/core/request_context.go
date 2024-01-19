package core

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type requestContextModel struct {
	TplHeader         string
	MetricsImportPath string
	ConfigsImportPath string
}

func newRequestContextModel(ctx *gencontext.GenContext) requestContextModel {
	// TODO: move paths to description
	commonPath := ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetCommonPackage()
	return requestContextModel{
		TplHeader:         ctx.GetWorkspace().TplHeader,
		MetricsImportPath: fmt.Sprintf("%s/metrics", commonPath),
		ConfigsImportPath: fmt.Sprintf("%s/configs", commonPath),
	}
}
