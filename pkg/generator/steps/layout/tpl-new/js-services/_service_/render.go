package service

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/components"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/pages"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed package.json.tpl
var packageJsonTemplate string

//go:embed nuxt.config.js.tpl
var nuxtConfigTemplate string

//go:embed dockerfile.tpl
var dockerfileTemplate string

func Render(ctx *gencontext.GenContext) error {
	packageJsonModel := newPackageJsonModel(ctx)
	packageJsonPath := ctx.GetWorkspace().GetJsServicePackageJsonAbsPath(ctx.GetServiceName())
	if err := render.RenderOrSkipTemplate(packageJsonTemplate, packageJsonModel, packageJsonPath); err != nil {
		return render.WrapError("package.json", err)
	}

	nuxtConfigModel := newNuxtConfigModel(ctx)
	nuxtConfigPath := ctx.GetWorkspace().GetJsServiceNuxtConfigAbsPath(ctx.GetServiceName())
	if err := render.RenderOrSkipTemplate(nuxtConfigTemplate, nuxtConfigModel, nuxtConfigPath); err != nil {
		return render.WrapError("nuxt.config.js", err)
	}

	dockerfileModel := newDockerfileModel(ctx)
	dockerfilePath := ctx.GetWorkspace().GetJsDockerfileAbsPath(ctx.GetServiceName())
	if err := render.RenderOrSkipTemplate(dockerfileTemplate, dockerfileModel, dockerfilePath); err != nil {
		return render.WrapError("Dockerfile.go", err)
	}

	if err := pages.Render(ctx); err != nil {
		return render.WrapError("pages", err)
	}

	if err := components.Render(ctx); err != nil {
		return render.WrapError("components", err)
	}

	return nil
}
