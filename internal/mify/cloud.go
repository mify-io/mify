package mify

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mify-io/mify/pkg/workspace"
)


func CloudInit(ctx *CliContext, basePath string) error {
	const CLOUD_URL = "https://cloud.mify.io"
	_, err := workspace.InitDescription(basePath)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	ctx.Logger.Printf("Please visit %s to receive token and paste it here:", CLOUD_URL)

	text, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read token from stdin: %w", err)
	}
	ctx.Config.APIToken = text
	err = SaveConfig(ctx.Config)
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	return nil
}

func CloudUpdateKubeconfig(ctx *CliContext, basePath string) error {
	return nil
}
