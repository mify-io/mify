package service

import gencontext "github.com/mify-io/mify/pkg/generator/gen-context"

type DockerfileModel struct {
	ServiceName string
}

func newDockerfileModel(ctx *gencontext.GenContext) DockerfileModel {
	return DockerfileModel{
		ServiceName: ctx.GetServiceName(),
	}
}
