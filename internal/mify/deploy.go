package mify

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mify-io/mify/internal/mify/util/docker"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

const (
	image = "mifyio/pipeline:latest"
)

func DeployMany(ctx *CliContext, deployEnv string, names []string) error {
	descr := ctx.MustGetWorkspaceDescription()

	if len(names) == 0 {
		names = descr.GetApiServices()
	}

	for _, name := range names {
		if err := deploy(ctx, deployEnv, name); err != nil {
			return fmt.Errorf("service '%s' deployment failed: %w", name, err)
		}
	}

	return nil
}

func deploy(ctx *CliContext, deployEnv string, serviceName string) error {
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

	cmdArgs := []string{"deploy", serviceName, "-p", "/repo"}
	if ctx.IsVerbose {
		cmdArgs = append(cmdArgs, "--verbose")
	}

	params := docker.DockerRunParams{
		Mounts: map[string]string{
			"/repo": ctx.WorkspacePath,
			// TODO: support other oses
			"/var/run/docker.sock": "/var/run/docker.sock",
		},
		Cmd: cmdArgs,
		Env: []string{
			"MIFY_API_TOKEN=" + strings.TrimSpace(ctx.Config.APIToken),
			"DEPLOY_ENVIRONMENT=" + deployEnv,
		},
		Tty: true,
	}

	err = docker.Run(ctx.GetCtx(), genContext.Logger, os.Stdout, image, params)
	if err != nil {
		return err
	}

	return nil
}
