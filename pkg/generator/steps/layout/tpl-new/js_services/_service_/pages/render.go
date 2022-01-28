package pages

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed index.vue
var indexTemplate string

func Render(ctx *gencontext.GenContext) error {
	indexModel := newIndexModel(ctx)
	indexPath := ctx.GetWorkspace().GetJsIndexAbsPath(ctx.GetServiceName())
	if err := render.RenderOrSkipTemplate(indexTemplate, indexModel, indexPath); err != nil {
		return render.WrapError("index.vue", err)
	}

	return nil
}
