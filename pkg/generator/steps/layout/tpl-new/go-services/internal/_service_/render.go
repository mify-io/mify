package service

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/_service_/app"
)

func Render(ctx *gencontext.GenContext) error {
	return app.Render(ctx)
}
