package mify

import (
	"os"
	"os/user"

	"github.com/mify-io/mify/internal/mify/util/docker"
	"go.uber.org/zap"
)

const (
	image = "mify-pipeline:latest"
)

func Deploy(ctx *CliContext, basePath string, args []string) error {
	if err := docker.Cleanup(ctx.GetCtx()); err != nil {
		return err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	// if err := docker.PullImage(ctx.GetCtx(), logger.Sugar(), os.Stdout, image); err != nil {
	// 	return err
	// }

	curUser, err := user.Current()
	if err != nil {
		return err
	}

	ctx.Logger.Println("running deploy")
	params := docker.DockerRunParams{
		User:   curUser,
		Mounts: map[string]string{"/repo": basePath},
		Cmd:    args,
	}

	err = docker.Run(ctx.GetCtx(), logger.Sugar(), os.Stdout, image, params)
	if err != nil {
		return err
	}

	return nil
}
