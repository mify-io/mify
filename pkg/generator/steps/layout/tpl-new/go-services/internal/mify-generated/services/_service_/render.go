package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/mify-generated/services/_service_/app"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/mify-generated/services/_service_/apputil"
)

func Render(ctx *gencontext.GenContext) error {
	if err := app.Render(ctx); err != nil {
		return err
	}

	if err := apputil.Render(ctx); err != nil {
		return err
	}

	return nil
}
