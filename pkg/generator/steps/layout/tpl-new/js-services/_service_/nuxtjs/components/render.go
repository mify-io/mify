package components

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/util/render"
)

//go:embed sample.vue.tpl
var sampleVueTemplate string

func Render(ctx *gencontext.GenContext) error {
	sampleVueModel := newSampleVueModel(ctx)
	sampleVuePath := ctx.GetWorkspace().GetJsSampleVueAbsPath(ctx.GetServiceName())
	if err := render.RenderOrSkipTemplate(sampleVueTemplate, sampleVueModel, sampleVuePath); err != nil {
		return render.WrapError("sample.vue", err)
	}

	return nil
}
