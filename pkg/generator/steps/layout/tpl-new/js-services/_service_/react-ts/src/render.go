package src

import (
	"embed"

	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	app "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/react-ts/src/app"
	generatedcore "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/react-ts/src/generated/core"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	if err := app.Render(ctx); err != nil {
		return err
	}
	if err := generatedcore.Render(ctx); err != nil {
		return err
	}
	dir := path.Join(ctx.GetWorkspace().GetJsServiceAbsPath(ctx.GetServiceName()), "src")

	return render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(dir, "reportWebVitals.ts")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "App.css")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "index.css")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "App.tsx")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "setupTests.ts")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "react-app-env.d.ts")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "App.test.tsx")).SetFlags(render.NewFlags().SkipExisting()),
		render.NewFile(ctx, path.Join(dir, "index.tsx")).SetFlags(render.NewFlags().SkipExisting()),
	)
}
