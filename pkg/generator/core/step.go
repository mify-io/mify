package core

import gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"

type Step interface {
	Name() string
	Execute(*gencontext.GenContext) error
}
