package services

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	service "github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/services/_service_"
	"github.com/mify-io/mify/pkg/util/render"
)

func Render(ctx *gencontext.GenContext) error {
	if err := service.Render(ctx); err != nil {
		return render.WrapError("service", err)
	}

	return nil
}
