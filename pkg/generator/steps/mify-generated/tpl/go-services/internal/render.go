package internal

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	mifygenerated "github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated"
)

func Render(ctx *gencontext.GenContext) error {
	if err := mifygenerated.Render(ctx); err != nil {
		return err
	}
	return nil
}
