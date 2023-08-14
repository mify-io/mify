package reactts

import (
	"embed"

	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/react-ts/public"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/react-ts/src"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext, dockerNginx bool) error {
	if err := src.Render(ctx); err != nil {
		return err
	}
	if err := public.Render(ctx); err != nil {
		return err
	}
	dockerfile := render.NewFile(ctx, ctx.GetWorkspace().GetJsDockerfileAbsPath(ctx.GetServiceName()))
	if dockerNginx {
		dockerfile = render.NewFile(ctx, ctx.GetWorkspace().GetJsDockerfileAbsPath(ctx.GetServiceName())).
			SetTemplateName("Dockerfile-nginx.tpl")
	}

	return render.RenderMany(
		templates,
		render.NewFile(ctx, ctx.GetWorkspace().GetJsServicePackageJsonAbsPath(ctx.GetServiceName())).
			SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(ctx.GetWorkspace().GetJsServiceAbsPath(ctx.GetServiceName()), "tsconfig.json")).
			SetFlags(render.NewFlags().SkipExisting()),
		dockerfile,
	)
}
