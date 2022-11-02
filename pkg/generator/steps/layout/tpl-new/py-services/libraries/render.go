package libraries

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"github.com/mify-io/mify/pkg/generator/steps/layout/tpl-new/py-services/libraries/generated"
)

func Render(ctx *gencontext.GenContext) error {
	return generated.Render(ctx)
}
