package common

import (
	"embed"
	"path"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/common/configs"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/common/logs"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/common/metrics"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed *.tpl
var templates embed.FS

func Render(ctx *gencontext.GenContext) error {
	basePath := path.Join(
		ctx.GetWorkspace().GetMifyGenerated(ctx.MustGetMifySchema()).GetCommonPath().Abs(),
	)
	if err := render.RenderMany(
		templates,
		render.NewFile(ctx, path.Join(basePath, "__init__.py")),
	); err != nil {
		return render.WrapError("__init__.py", err)
	}

	if err := configs.Render(ctx); err != nil {
		return err
	}
	if err := logs.Render(ctx); err != nil {
		return err
	}
	if err := metrics.Render(ctx); err != nil {
		return err
	}
	return nil
}
