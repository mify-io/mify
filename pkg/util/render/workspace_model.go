package render

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

type GoServiceModel struct {
	Name string
}

type WorkspaceModel struct {
	Name       string
	MifyGeneratedCommonPackage string
	MifyGeneratedServicePackage string
}

func NewWorkspaceModel(context *gencontext.GenContext) *WorkspaceModel {
	mifyGen := context.GetWorkspace().GetMifyGenerated(context.MustGetMifySchema())
	return &WorkspaceModel{
		Name:       context.GetWorkspace().Name,
		MifyGeneratedCommonPackage: mifyGen.GetCommonPackage(),
		MifyGeneratedServicePackage: mifyGen.GetServicePackage(),
	}
}

