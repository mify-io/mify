package configs

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type dynamicModel struct {
}

func newDynamicModel(ctx *gencontext.GenContext) dynamicModel {
	return dynamicModel{}
}
