package generated

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/_service_/generated/app"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/_service_/generated/core"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/_service_/generated/openapi"
)

func Render(ctx *gencontext.GenContext) error {
	if err := core.Render(ctx); err != nil {
		return err
	}

	if err := app.Render(ctx); err != nil {
		return err
	}

	if err := openapi.Render(ctx); err != nil {
		return err
	}

	return nil
}
