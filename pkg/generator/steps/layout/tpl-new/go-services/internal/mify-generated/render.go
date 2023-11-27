package mifygenerated

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	services "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/mify-generated/services"
)

func Render(ctx *gencontext.GenContext) error {
	if err := services.Render(ctx); err != nil {
		return err
	}
	return nil
}
