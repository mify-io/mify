package generated

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/pkg/generated/configs"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/pkg/generated/consul"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/pkg/generated/logs"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/pkg/generated/metrics"
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
