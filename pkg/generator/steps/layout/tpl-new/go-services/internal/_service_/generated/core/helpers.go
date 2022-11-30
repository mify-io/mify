package core

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type helpersModel struct {
	TplHeader string
}

func newHelpersModel(ctx *gencontext.GenContext) helpersModel {
	return helpersModel{
		TplHeader: ctx.GetWorkspace().TplHeader,
	}
}
