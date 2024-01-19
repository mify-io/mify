package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/services/_service_/app"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/services/_service_/core"
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
