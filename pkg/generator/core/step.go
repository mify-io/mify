package core

import "github.com/chebykinn/mify/pkg/generator/context"

type Step interface {
	Name() string
	Execute(*context.GenContext) error
}
