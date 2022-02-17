package generated

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_/generated/app"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_/generated/core"
)

func Render(ctx *gencontext.GenContext) error {
	if err := core.Render(ctx); err != nil {
		return err
	}

	if err := app.Render(ctx); err != nil {
		return err
	}

	return nil
}
