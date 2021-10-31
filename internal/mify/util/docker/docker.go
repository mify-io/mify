package docker

import (
	"context"
	"fmt"
	"os"
	"os/user"

	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerRunParams struct {
	User *user.User
	Mounts map[string]string
	Cmd []string
}

func Run(ctx context.Context, image string, params DockerRunParams) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	spinner := util.CreateWaitSpinner()

	spinner.Suffix = " docker: pulling "+image
	spinner.Start()
	_, err = cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	spinner.Stop()
	// io.Copy(os.Stdout, reader)

	m := make([]mount.Mount, 0, len(params.Mounts))
	for target, src := range params.Mounts {
		m = append(m, mount.Mount{
			Type: mount.TypeBind,
			Source: src,
			Target: target,
		})
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		User: params.User.Uid+":"+params.User.Gid,
		Image: image,
		Cmd: params.Cmd,
		Tty: false,
	}, &container.HostConfig{Mounts: m}, nil, nil, "")
	if err != nil {
		return err
	}
	defer func() {
		removeOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		if err := cli.ContainerRemove(ctx, resp.ID, removeOptions); err != nil {
			fmt.Printf("unable to remove container: %s", err)
		}
	}()

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	spinner.Suffix = " docker: running "+image
	spinner.Start()
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
	spinner.Stop()

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if exitCode != 0 {
		return fmt.Errorf("process exited with code: %d", exitCode)
	}

	return nil
}
