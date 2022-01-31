package openapi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	"gopkg.in/yaml.v2"
)

type clientsDiff struct {
	schemaChanged map[string]struct{}
	added         map[string]struct{}
	removed       map[string]struct{}
}

func (c clientsDiff) Empty() bool {
	return len(c.schemaChanged) == 0 &&
		len(c.added) == 0 &&
		len(c.removed) == 0
}

func calcClientsDiff(ctx *gencontext.GenContext, openAPIGenerator *OpenAPIGenerator) (clientsDiff, error) {
	wrapError := func(err error) error {
		return fmt.Errorf("can't calculate clients diff: %w", err)
	}

	schemaChanged := make(map[string]struct{})
	added := make(map[string]struct{})
	removed := make(map[string]struct{})

	currentClients := getCurrentClients(ctx)

	oldClients, err := getOldClients(ctx)
	if err != nil {
		return clientsDiff{}, wrapError(err)
	}

	if oldClients == nil {
		for cl := range currentClients {
			added[cl] = struct{}{}
		}
		return clientsDiff{added: added}, nil
	}

	for cl := range currentClients {
		if _, ok := oldClients[cl]; !ok {
			added[cl] = struct{}{}
			continue
		}

		schemaDirPath := ctx.GetWorkspace().GetApiSchemaDirRelPath(cl)
		needGenerate, err := openAPIGenerator.NeedGenerateClient(ctx, schemaDirPath)
		if err != nil {
			return clientsDiff{}, wrapError(err)
		}

		if needGenerate {
			schemaChanged[cl] = struct{}{}
		}
	}

	for cl := range oldClients {
		if _, ok := currentClients[cl]; ok {
			continue
		}
		removed[cl] = struct{}{}
	}

	return clientsDiff{schemaChanged: schemaChanged, added: added, removed: removed}, nil
}

func getCurrentClients(ctx *gencontext.GenContext) map[string]struct{} {
	currentClients := ctx.MustGetMifySchema().OpenAPI.Clients
	currentClientsSet := make(map[string]struct{})
	for cl := range currentClients {
		currentClientsSet[cl] = struct{}{}
	}

	return currentClientsSet
}

// Loads old clients set from cache dir
func getOldClients(ctx *gencontext.GenContext) (map[string]struct{}, error) {
	tmpDir := ctx.GetWorkspace().GetServiceCacheDirectory(ctx.GetServiceName())
	f, err := os.Open(filepath.Join(tmpDir, CLIENTS_FILENAME))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var oldClients []string
	err = yaml.NewDecoder(f).Decode(&oldClients)
	if err != nil {
		return nil, err
	}

	oldClientsSet := make(map[string]struct{})
	for _, cl := range oldClients {
		oldClientsSet[cl] = struct{}{}
	}

	return oldClientsSet, nil
}
