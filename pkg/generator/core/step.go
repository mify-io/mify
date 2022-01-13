package core

import gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"

type StepResult int

const (
	Done      StepResult = 0
	RepeatAll StepResult = 1
)

type Step interface {
	Name() string
	Execute(*gencontext.GenContext) (StepResult, error)
}
