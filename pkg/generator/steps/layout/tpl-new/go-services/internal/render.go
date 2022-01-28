package internal

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	service "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_"
)

func Render(ctx *gencontext.GenContext) error {
	return service.Render(ctx)
}
