package gotpl

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/postgres/tpl/go/config"
)

func Render(ctx *gencontext.GenContext) error {
	return config.Render(ctx)
}
