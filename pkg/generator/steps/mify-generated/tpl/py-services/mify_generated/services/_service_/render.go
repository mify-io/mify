package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/services/_service_/app"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/services/_service_/core"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/services/_service_/openapi"
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
