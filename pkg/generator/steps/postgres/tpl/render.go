package tpl

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	gotpl "github.com/mify-io/mify/pkg/generator/steps/postgres/tpl/go"
)

func RenderGo(ctx *gencontext.GenContext) error {
	if err := gotpl.Render(ctx); err != nil {
		return err
	}

	return nil
}
