package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/chebykinn/mify/internal/mify/config"
	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/chebykinn/mify/internal/mify/service/generate"
	"github.com/chebykinn/mify/internal/mify/util"
	"github.com/chebykinn/mify/pkg/mifyconfig"
	"gopkg.in/yaml.v2"
)

const (
	CLIENTS_FILENAME    = ".clients.yaml"
)

type clientsDiff struct {
	added map[string]struct{}
	removed map[string]struct{}
}

func makeClientsContext(conf mifyconfig.ServiceConfig, basePath string) ([]OpenAPIClientContext, error) {
	clientsCtxList := make([]OpenAPIClientContext, 0, len(conf.OpenAPI.Clients))
	for clientName := range conf.OpenAPI.Clients {
		clientSchemaPath := fmt.Sprintf(apiSchemaPath, clientName)
		if _, err := os.Stat(filepath.Join(basePath, clientSchemaPath)); errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("schema not found while generating client for: %s", clientName)
		}

		packageName := generate.MakePackageName(clientName)
		fieldName := generate.SnakeCaseToCamelCase(generate.SanitizeServiceName(clientName), false)
		methodName := generate.SnakeCaseToCamelCase(generate.SanitizeServiceName(clientName), true)
		clientsCtxList = append(clientsCtxList, OpenAPIClientContext{
			ClientName: clientName,
			PackageName: packageName,
			PrivateFieldName: fieldName,
			PublicMethodName: methodName,
		})

	}
	sort.Slice(clientsCtxList, func(i, j int) bool {
		return clientsCtxList[i].ClientName < clientsCtxList[j].ClientName
	})
	return clientsCtxList, nil
}

func generateClientsContextStep(ctx *core.Context, pool *util.JobPool, serviceCtx Context, conf mifyconfig.ServiceConfig) (clientsDiff, error) {
	clientsCtx, err := makeClientsContext(conf, serviceCtx.Workspace.BasePath)
	if err != nil {
		return clientsDiff{}, err
	}
	serviceCtx.OpenAPI.Clients = clientsCtx

	list := make([]string, 0, len(clientsCtx))
	for _, c := range clientsCtx {
		list = append(list, c.ClientName)
	}

	tmpDir := config.GetServiceCacheDirectory(serviceCtx.Workspace.BasePath, serviceCtx.ServiceName)
	diff, err := getClientsDiff(ctx, serviceCtx.Workspace.BasePath, tmpDir, list)
	if err != nil {
		return clientsDiff{}, err
	}
	if len(diff.added) == 0 && len(diff.removed) == 0 {
		return clientsDiff{}, ErrSkip
	}
	pool.AddJob(util.Job{
		Name: "generate:clients-context",
		Func: func(ctx *core.Context) error {
			err := generateClientsContext(ctx, serviceCtx, conf, list)
			if err != nil {
				return err
			}
			return nil
		},
	})
	if err := pool.Run(); err != nil {
		return clientsDiff{}, err
	}
	return diff, nil
}

func generateClientsContext(ctx *core.Context, serviceCtx Context, conf mifyconfig.ServiceConfig, list []string) error {
	tmpDir := config.GetServiceCacheDirectory(serviceCtx.Workspace.BasePath, serviceCtx.ServiceName)
	switch serviceCtx.Language {
		case mifyconfig.ServiceLanguageGo:
			subPath := mifyconfig.GoServicesRoot+"/internal/#svc#/generated/core"
			// FIXME: regenerate only clients file
			if err := RenderTemplateTreeSubPath(ctx, serviceCtx, subPath); err != nil {
				return err
			}
		case mifyconfig.ServiceLanguageJs:
			subPath := mifyconfig.JsServicesRoot+"/#svc#/generated/core"
			if err := RenderTemplateTreeSubPath(ctx, serviceCtx, subPath); err != nil {
				return err
			}
	}

	err := updateClientsList(ctx, serviceCtx.Workspace.BasePath, tmpDir, list)
	if err != nil {
		return err
	}

	return nil
}

func updateClientsList(ctx *core.Context, basePath string, tmpDir string, list []string) error {
	ctx.Logger.Printf("updating clients list in: %s", tmpDir)

	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to write clients list file: %w", err)
	}

	f, err := os.OpenFile(filepath.Join(tmpDir, CLIENTS_FILENAME), os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to write clients list file: %w", err)
	}

	err = yaml.NewEncoder(f).Encode(list)
	if err != nil {
		return fmt.Errorf("failed to write clients list file: %w", err)
	}

	return nil
}

func getClientsDiff(ctx *core.Context, basePath string, tmpDir string, newList []string) (clientsDiff, error) {
	f, err := os.Open(filepath.Join(tmpDir, CLIENTS_FILENAME))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			mp := map[string]struct{}{}
			for _, cl := range newList {
				mp[cl] = struct{}{}
			}
			return clientsDiff{added: mp}, nil
		}
		return clientsDiff{}, fmt.Errorf("failed to compare clients: %w", err)
	}

	var oldList []string
	err = yaml.NewDecoder(f).Decode(&oldList)
	if err != nil {
		return clientsDiff{}, fmt.Errorf("failed to compare clients: %w", err)
	}

	oldSet := map[string]struct{}{}
	for _, cl := range oldList {
		oldSet[cl] = struct{}{}
	}

	newClients := map[string]struct{}{}
	for _, cl := range newList {
		if _, ok := oldSet[cl]; ok {
			continue
		}
		newClients[cl] = struct{}{}
	}

	newSet := map[string]struct{}{}
	for _, cl := range newList {
		newSet[cl] = struct{}{}
	}

	delClients := map[string]struct{}{}
	for _, cl := range oldList {
		if _, ok := newSet[cl]; ok {
			continue
		}
		delClients[cl] = struct{}{}
	}

	return clientsDiff{ added: newClients, removed: delClients}, nil
}
