package goservices

import (
	_ "embed"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal"
	"github.com/mify-io/mify/pkg/util/render"
)

func Render(ctx *gencontext.GenContext) error {
	if err := internal.Render(ctx); err != nil {
		return render.WrapError("app", err)
	}
	return nil
}
