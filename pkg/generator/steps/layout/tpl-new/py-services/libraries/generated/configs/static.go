package configs

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type staticModel struct {
}

func newStaticModel(ctx *gencontext.GenContext) staticModel {
	return staticModel{}
}
