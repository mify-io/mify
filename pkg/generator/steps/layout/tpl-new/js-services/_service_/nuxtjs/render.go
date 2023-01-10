package nuxtjs

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/nuxtjs/components"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/nuxtjs/generated/core"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/nuxtjs/pages"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/nuxtjs/plugins"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/nuxtjs/store"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed package.json.tpl
var packageJsonTemplate string

//go:embed yarn.lock.tpl
var yarnLockTemplate string

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

	yarnLockModel := newYarnLockModel(ctx)
	yarnLockPath := ctx.GetWorkspace().GetJsServiceYarnLockAbsPath(ctx.GetServiceName())
	if err := render.RenderOrSkipTemplate(yarnLockTemplate, yarnLockModel, yarnLockPath); err != nil {
		return render.WrapError("yarn.lock", err)
	}

	if err := generatedcore.Render(ctx); err != nil {
		return render.WrapError("generated", err)
	}

	if err := pages.Render(ctx); err != nil {
		return render.WrapError("pages", err)
	}

	if err := components.Render(ctx); err != nil {
		return render.WrapError("components", err)
	}

	if err := store.Render(ctx); err != nil {
		return render.WrapError("store", err)
	}

	if err := plugins.Render(ctx); err != nil {
		return render.WrapError("plugins", err)
	}

	return nil
}
