package internal

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	service "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/pkg"
)

func Render(ctx *gencontext.GenContext) error {
	if err := service.Render(ctx); err != nil {
		return err
	}
	if err := pkg.Render(ctx); err != nil {
		return err
	}
	return nil
}
