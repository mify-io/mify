package tplnew

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	goservices "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services"
	jsservices "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services"
)

func RenderGo(ctx *gencontext.GenContext) error {
	if err := goservices.Render(ctx); err != nil {
		return err
	}

	return nil
}

func RenderJs(ctx *gencontext.GenContext) error {
	if err := jsservices.Render(ctx); err != nil {
		return err
	}

	return nil
}
