package mifygenerated

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/common"
	services "github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/services"
)

func Render(ctx *gencontext.GenContext) error {
	if err := services.Render(ctx); err != nil {
		return err
	}
	if err := common.Render(ctx); err != nil {
		return err
	}
	return nil
}
