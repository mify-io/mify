package docker

import (
	"context"
	"fmt"
	"io"
	"os"
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

func PullImage(ctx context.Context, logger *zap.SugaredLogger, image string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	logger.Infof("pulling image: %s", image)
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	// TODO: do smth else
	io.Copy(os.Stdout, reader)

	return nil
}

func Run(ctx context.Context, logger *zap.SugaredLogger, image string, params DockerRunParams) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	if params.PullImage {
		if err := PullImage(ctx, logger, image); err != nil {
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

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		User:  params.User.Uid + ":" + params.User.Gid,
		Image: image,
		Cmd:   params.Cmd,
		Tty:   false,
		Labels: map[string]string{
			mifyContainerLabel: "",
		},
	}, &container.HostConfig{Mounts: m}, nil, nil, "")
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
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case st := <-statusCh:
		exitCode = st.StatusCode
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return err
	}

	// TODO: do smth else
	io.Copy(os.Stdout, out)
	if exitCode != 0 {
		return fmt.Errorf("process exited with code: %d", exitCode)
	}

	return nil
}
