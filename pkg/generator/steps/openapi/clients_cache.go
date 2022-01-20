package openapi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/internal/mify/config"
	gencontext "github.com/chebykinn/mify/pkg/generator/gen-context"
	"gopkg.in/yaml.v2"
)

const (
	CLIENTS_FILENAME = ".clients.yaml"
)

func updateClientsList(ctx *gencontext.GenContext) error {
	tmpDir := config.GetServiceCacheDirectory(ctx.GetWorkspace().BasePath, ctx.GetServiceName())

	ctx.Logger.Infof("updating clients list in: %s", tmpDir)

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to write clients list file: %w", err)
	}

	f, err := os.OpenFile(filepath.Join(tmpDir, CLIENTS_FILENAME), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to write clients list file: %w", err)
	}

	list := make([]string, 0, len(ctx.MustGetMifySchema().OpenAPI.Clients))
	for client := range ctx.MustGetMifySchema().OpenAPI.Clients {
		list = append(list, client)
	}

	err = yaml.NewEncoder(f).Encode(list)
	if err != nil {
		return fmt.Errorf("failed to write clients list file: %w", err)
	}

	return nil
}
