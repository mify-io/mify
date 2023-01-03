package core

import (
	"fmt"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type helpersModel struct {
	TplHeader         string
	ConfigsImportPath string
}

func newHelpersModel(ctx *gencontext.GenContext) helpersModel {
	return helpersModel{
		TplHeader:         ctx.GetWorkspace().TplHeader,
		ConfigsImportPath: fmt.Sprintf("%s/internal/pkg/generated/configs", ctx.GetWorkspace().GetGoModule()),
	}
}
