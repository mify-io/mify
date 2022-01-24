package app

import (
	"{{.CoreInclude}}"
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
