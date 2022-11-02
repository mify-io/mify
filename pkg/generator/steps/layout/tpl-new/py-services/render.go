package pyservices

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	service "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/_service_"
	libraries "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/libraries"
)

func Render(ctx *gencontext.GenContext) error {
	if err := service.Render(ctx); err != nil {
		return err
	}
	if err := libraries.Render(ctx); err != nil {
		return err
	}
	return nil
}
