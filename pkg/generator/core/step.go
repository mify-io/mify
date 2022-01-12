package core

import "context"

type ExecuteFunc func(*context.Context) *context.Context

type Step interface {
	Name() string
	ExecuteFunc() ExecuteFunc
}
