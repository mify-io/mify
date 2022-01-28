package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type dockerfileModel struct {
	ServiceName string
}

func newDockerfileModel(ctx *gencontext.GenContext) dockerfileModel {
	return dockerfileModel{
		ServiceName: ctx.GetServiceName(),
	}
}
