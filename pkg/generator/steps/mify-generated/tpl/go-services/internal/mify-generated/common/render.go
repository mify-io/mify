package common

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/common/configs"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/common/consul"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/common/logs"
	"github.com/mify-io/mify/pkg/generator/steps/mify-generated/tpl/go-services/internal/mify-generated/common/metrics"
)

func Render(ctx *gencontext.GenContext) error {
	if err := configs.Render(ctx); err != nil {
		return err
	}
	if err := consul.Render(ctx); err != nil {
		return err
	}
	if err := logs.Render(ctx); err != nil {
		return err
	}
	if err := metrics.Render(ctx); err != nil {
		return err
	}
	return nil
}
