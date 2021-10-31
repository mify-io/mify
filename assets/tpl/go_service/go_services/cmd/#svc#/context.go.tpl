package main

import (
	"{{.GoModule}}/internal/{{.ServiceName}}/misc"
)

{{- $ctxStructName := svcUserCtxName .}}

type {{$ctxStructName}} struct {
	// Append your dependencies here
}

func New{{$ctxStructName}} (mifyServiceContext misc.MifyServiceContext) ({{$ctxStructName}}, error) {
	context := Service1Context{
		// Here you can initialize your dependencies
	}
	return context, nil
}
