package mify

import (
	"io"
	"os"

	"github.com/mify-io/mify/internal/mify/util/docker"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

const (
	image = "mifyio/pipeline:latest"
)

func Deploy(ctx *CliContext, deployEnv string, serviceName string) error {
	ctx.Logger.Printf("Deploying service %s to %s, environment: %s", serviceName,
		ctx.workspaceDescription.Config.WorkspaceName, deployEnv)
	err := ServiceGenerate(ctx, ctx.WorkspacePath, serviceName)
	if err != nil {
		return err
	}
	// TODO: maybe separate logger
	genContext := gencontext.NewGenContext(ctx.Ctx, serviceName, *ctx.workspaceDescription)

	if err := docker.Cleanup(ctx.GetCtx()); err != nil {
		return err
	}

	if err := docker.PullImage(ctx.GetCtx(), genContext.Logger, io.Discard, image); err != nil {
		return err
	}

	params := docker.DockerRunParams{
		Mounts: map[string]string{
			"/repo": ctx.WorkspacePath,
			// TODO: support other oses
			"/var/run/docker.sock": "/var/run/docker.sock",
		},
		Cmd: []string{"deploy", serviceName, "-p", "/repo"},
		Env: []string{
			"MIFY_API_TOKEN=" + ctx.Config.APIToken,
			"DEPLOY_ENVIRONMENT=" + deployEnv,
		},
	}

	err = docker.Run(ctx.GetCtx(), genContext.Logger, os.Stdout, image, params)
	if err != nil {
		return err
	}

	return nil
}
