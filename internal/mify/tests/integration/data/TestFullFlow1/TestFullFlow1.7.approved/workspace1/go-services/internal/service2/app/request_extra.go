package app

import (
	"example.com/namespace/workspace1/go-services/internal/service2/generated/core"
)

type RequestExtra struct {
	// Append your dependencies here
}

func NewRequestExtra(ctx *core.MifyServiceContext) (*RequestExtra, error) {
	// Here you can do your custom service initialization, prepare dependencies
	extra := &RequestExtra{
		// Here you can initialize your dependencies
	}
	return extra, nil
}
