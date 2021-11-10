package core

import (
	"context"
	"log"
	"os"
)

type Context struct {
	Logger *log.Logger
	Ctx context.Context
	Cancel context.CancelFunc
}

func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		Logger: log.New(os.Stdout, "", 0),
		Ctx: ctx,
		Cancel: cancel,
	}
}
