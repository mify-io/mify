package tpl

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	goservices "github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services"
	// jsservices "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/js-services"
	// pyservices "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services"
)

func RenderGo(ctx *gencontext.GenContext) error {
	if err := goservices.Render(ctx); err != nil {
		return err
	}

	return nil
}

func RenderJs(ctx *gencontext.GenContext) error {
	// if err := jsservices.Render(ctx); err != nil {
		// return err
	// }

	return nil
}

func RenderPy(ctx *gencontext.GenContext) error {
	// if err := pyservices.Render(ctx); err != nil {
		// return err
	// }

	return nil
}
