package pyservices

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	mifygenerated "github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated"
	"github.com/mify-io/mify/pkg/util/render"
)

func Render(ctx *gencontext.GenContext) error {
	if err := mifygenerated.Render(ctx); err != nil {
		return render.WrapError("mifygenerated", err)
	}

	return nil
}
