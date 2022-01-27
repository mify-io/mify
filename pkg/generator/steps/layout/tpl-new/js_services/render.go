package jsservices

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed package.json.tpl
var packageJsonTemplate string

func Render(ctx *gencontext.GenContext) error {
	packageJsonModel := NewPackageJsonModel(ctx)
	packageJsonPath := ctx.GetWorkspace().GetJsPackageJsonAbsPath()
	if err := render.RenderTemplate(packageJsonTemplate, packageJsonModel, packageJsonPath); err != nil {
		return render.WrapError("package.json", err)
	}

	return nil
}
