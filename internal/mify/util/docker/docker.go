package docker

import (
	"context"
	"fmt"
	"io"
	"os/user"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

const (
	mifyContainerLabel = "mify.container"
)

type DockerRunParams struct {
	User      *user.User
	Mounts    map[string]string
	Cmd       []string
	PullImage bool
	Env       []string
	Tty       bool
}

func removeContainer(ctx context.Context, client *client.Client, id string) error {
	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	return client.ContainerRemove(ctx, id, removeOptions)
}

func Cleanup(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", mifyContainerLabel),
			filters.Arg("status", "created"),
			filters.Arg("status", "restarting"),
			filters.Arg("status", "running"),
			filters.Arg("status", "removing"),
			filters.Arg("status", "paused"),
			filters.Arg("status", "exited"),
			filters.Arg("status", "dead"),
		),
	})
	if err != nil {
		return err
	}
	for _, container := range containers {
		if err := removeContainer(ctx, cli, container.ID); err != nil {
			return err
		}
	}
	return nil
}

func PullImage(ctx context.Context, logger *zap.SugaredLogger, dockerLogs io.Writer, image string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	logger.Infof("pulling image: %s", image)
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	_, err = io.Copy(dockerLogs, reader)
	if err != nil {
		return err
	}

	return nil
}

func Run(
	ctx context.Context,
	logger *zap.SugaredLogger,
	dockerLogs io.Writer,
	image string,
	params DockerRunParams) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	if params.PullImage {
		if err := PullImage(ctx, logger, dockerLogs, image); err != nil {
			return err
		}
	}

	m := make([]mount.Mount, 0, len(params.Mounts))
	for target, src := range params.Mounts {
		m = append(m, mount.Mount{
			Type:   mount.TypeBind,
			Source: src,
			Target: target,
		})
	}
	conf := &container.Config{
		Image: image,
		Cmd:   params.Cmd,
		Tty:   params.Tty,
		Labels: map[string]string{
			mifyContainerLabel: "",
		},
		Env: params.Env,
	}
	if params.User != nil {
		conf.User = params.User.Uid + ":" + params.User.Gid
	}

	resp, err := cli.ContainerCreate(ctx, conf, &container.HostConfig{Mounts: m}, nil, nil, "")
	if err != nil {
		return err
	}
	defer func() {
		if err := removeContainer(ctx, cli, resp.ID); err != nil {
			logger.Infof("unable to remove container: %s", err)
		}
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	logger.Infof("running image: %s", image)
	var exitCode int64
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	go func() {
		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		if err != nil {
			logger.Errorf("Failed to read container logs: %s", err)
			return
		}

		_, err = io.Copy(dockerLogs, out)
		if err != nil {
			logger.Errorf("Failed to read container logs: %s", err)
			return
		}
	}()
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case st := <-statusCh:
		exitCode = st.StatusCode
	}

	if exitCode != 0 {
		return fmt.Errorf("process exited with code: %d", exitCode)
	}

	return nil
}
