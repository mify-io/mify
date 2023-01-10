package jsservices

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	nuxtjs "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/nuxtjs"
	reactts "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services/_service_/react-ts"
	"github.com/mify-io/mify/pkg/util/render"
)

func Render(ctx *gencontext.GenContext) error {
	if ctx.GetMifySchema().Template == "nuxtjs" {
		if err := nuxtjs.Render(ctx); err != nil {
			return render.WrapError("nuxtjs", err)
		}
		return nil
	}
	if ctx.GetMifySchema().Template == "react-ts" {
		if err := reactts.Render(ctx); err != nil {
			return render.WrapError("reactts", err)
		}
		return nil
	}

	return nil
}
