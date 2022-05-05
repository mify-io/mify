package prepare

import (
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

func execute(ctx *gencontext.GenContext) error {
	if err := preparePython(ctx); err != nil {
		return err
	}
	return nil
}
