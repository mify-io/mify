package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_/app"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_/generated"
	"github.com/mify-io/mify/pkg/util/render"
)

func Render(ctx *gencontext.GenContext) error {
	if err := generated.Render(ctx); err != nil {
		return render.WrapError("generated", err)
	}

	if err := app.Render(ctx); err != nil {
		return render.WrapError("app", err)
	}

	return nil
}
