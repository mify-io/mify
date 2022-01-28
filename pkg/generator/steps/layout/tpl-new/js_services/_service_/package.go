package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type packageJsonModel struct {
	ServiceName string
}

func newPackageJsonModel(ctx *gencontext.GenContext) packageJsonModel {
	return packageJsonModel{
		ctx.GetServiceName(),
	}
}
