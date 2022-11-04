package golang

import (
	"github.com/mify-io/mify/pkg/generator/core"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

func Execute(ctx *gencontext.GenContext) (core.StepResult, error) {

	if err := execute(ctx); err != nil {
		return core.Done, err
	}

	return core.Done, nil
}
