package mify

import (
	"context"

	"github.com/mify-io/mify/internal/mify/util/docker"
	"github.com/mify-io/mify/pkg/generator"
	"github.com/mify-io/mify/pkg/workspace"
)

func ServiceGenerate(ctx *CliContext, basePath string, name string) error {
	descr, err := workspace.InitDescription(basePath)
	if err != nil {
		return err
	}

	genPipeline := generator.BuildServicePipeline()
	if err := genPipeline.Execute(ctx.Ctx, name, descr); err != nil {
		return err
	}

	return nil
}

func Cleanup() error {
	err := docker.Cleanup(context.Background())
	if err != nil {
		return err
	}

	return nil
}
