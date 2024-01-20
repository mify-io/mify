package mifygenerated

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/common"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/py-services/mify_generated/services"
	"github.com/mify-io/mify/pkg/util/render"
)

func Render(ctx *gencontext.GenContext) error {
	if err := common.Render(ctx); err != nil {
		return render.WrapError("common", err)
	}
	if err := services.Render(ctx); err != nil {
		return render.WrapError("services", err)
	}

	return nil
}
