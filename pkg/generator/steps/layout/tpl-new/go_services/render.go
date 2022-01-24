package goservices

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go_services/internal"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed go.mod.tpl
var goModTemplate string

//go:embed go.sum.tpl
var goSumTemplate string

func Render(ctx *gencontext.GenContext) error {
	goModModel := NewGoModModel(ctx.GetWorkspace().GetGoModule())
	goModPath := ctx.GetWorkspace().GetGoModAbsPath()
	if err := render.RenderOrSkipTemplate(goModTemplate, goModModel, goModPath); err != nil {
		return render.WrapError("go.mod", err)
	}

	goSumModel := NewGoSumModel()
	goSumPath := ctx.GetWorkspace().GetGoSumAbsPath()
	if err := render.RenderOrSkipTemplate(goSumTemplate, goSumModel, goSumPath); err != nil {
		return render.WrapError("go.sum", err)
	}

	if err := internal.Render(ctx); err != nil {
		return render.WrapError("app", err)
	}

	return nil
}
