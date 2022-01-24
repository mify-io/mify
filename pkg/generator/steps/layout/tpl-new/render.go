package tplnew

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	goservices "github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go_services"
)

func RenderGo(ctx *gencontext.GenContext) error {
	return goservices.Render(ctx)
}
