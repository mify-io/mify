package jsservices

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	service "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js_services/_service_"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed package.json.tpl
var packageJsonTemplate string

func Render(ctx *gencontext.GenContext) error {
	packageJsonModel := newPackageJsonModel(ctx)
	packageJsonPath := ctx.GetWorkspace().GetJsPackageJsonAbsPath()
	if err := render.RenderTemplate(packageJsonTemplate, packageJsonModel, packageJsonPath); err != nil {
		return render.WrapError("package.json", err)
	}

	if err := service.Render(ctx); err != nil {
		return render.WrapError("service", err)
	}

	return nil
}
