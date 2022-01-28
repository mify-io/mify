package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type PackageJsonModel struct {
	ServiceName string
}

func newPackageJsonModel(ctx *gencontext.GenContext) PackageJsonModel {
	return PackageJsonModel{
		ctx.GetServiceName(),
	}
}
