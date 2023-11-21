package pkg

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/go-services/internal/pkg/generated"
)

func Render(ctx *gencontext.GenContext) error {
	return generated.Render(ctx)
}
