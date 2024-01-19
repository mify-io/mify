package migrate

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
	migrator "github.com/mify-io/mify/pkg/generator/migrate"
	"github.com/mify-io/mify/pkg/mifyconfig"
)

func execute(ctx *gencontext.GenContext) error {
	mifySchema := ctx.MustGetMifySchema()
	if !mifySchema.Components.Layout.Enabled {
		ctx.Logger.Info("not migrating service without layout enabled")
		return nil
	}

	switch mifySchema.Language {
	case mifyconfig.ServiceLanguageGo:
		if err := migrateGo(ctx); err != nil {
			return fmt.Errorf("can't migrate go template: %w", err)
		}
	case mifyconfig.ServiceLanguageJs:
		if err := migrateJs(ctx); err != nil {
			return fmt.Errorf("can't migrate js template: %w", err)
		}
	case mifyconfig.ServiceLanguagePython:
		if err := migratePython(ctx); err != nil {
			return fmt.Errorf("can't migrate python template: %w", err)
		}
	}
	return nil
}

func migrateGo(ctx *gencontext.GenContext) error {
	toRemove := []string {
		"internal/pkg/generated/configs",
		"internal/pkg/generated/consul",
		"internal/pkg/generated/logs",
		"internal/pkg/generated/metrics",
		fmt.Sprintf("internal/%s/generated/api", ctx.GetServiceName()),
		fmt.Sprintf("internal/%s/generated/app", ctx.GetServiceName()),
		fmt.Sprintf("internal/%s/generated/apputil", ctx.GetServiceName()),
		fmt.Sprintf("internal/%s/generated/core", ctx.GetServiceName()),
		fmt.Sprintf("internal/%s/generated/postgres", ctx.GetServiceName()),
		fmt.Sprintf("internal/%s/generated/.openapi-generator", ctx.GetServiceName()),
		fmt.Sprintf("internal/%s/generated/.openapi-generator-ignore", ctx.GetServiceName()),
	}
	toRemoveDirs := []string {
		"internal/pkg/generated",
		"internal/pkg",
		fmt.Sprintf("internal/%s/generated", ctx.GetServiceName()),
	}
	toReplaceLookupPaths := []string {
		fmt.Sprintf("internal/%s", ctx.GetServiceName()),
		fmt.Sprintf("cmd/%s", ctx.GetServiceName()),
	}

	checkIfNeededPath := path.Join(
		ctx.GetWorkspace().BasePath,
		mifyconfig.GoServicesRoot,
		"cmd",
		ctx.GetServiceName(),
		"main.go",
	)
	prevSvcPath := fmt.Sprintf("%s/internal/%s/generated", ctx.GetWorkspace().GetGoModule(), ctx.GetServiceName())
	ok, err := migrator.FileContainsString(checkIfNeededPath, prevSvcPath)
	if err != nil {
		return err
	}
	if !ok {
		ctx.Logger.Infof("not doing old generated files migration")
		return nil
	}
	var checkVcsFunc = func(f string) error {
		filePath := path.Join(mifyconfig.GoServicesRoot, f)
		err := filepath.WalkDir(filePath, func(p string, d fs.DirEntry, ferr error) error {
			if d == nil || d.IsDir() {
				return nil
			}
			hasChanged, err := ctx.GetVcsIntegration().FileHasUncommitedChanges(p)
			if err != nil {
				return err
			}
			if hasChanged {
				return fmt.Errorf("migration can't be applied, path %s has uncommitted changes", filePath)
			}
			return nil
		})
		return err
	}
	for _, f := range toReplaceLookupPaths {
		if err := checkVcsFunc(f); err != nil {
			return err
		}
	}
	for _, f := range toRemove {
		filePath := path.Join(ctx.GetWorkspace().BasePath, mifyconfig.GoServicesRoot, f)
		ctx.Logger.Debugf("removing %s", filePath)
		if err := os.RemoveAll(filePath); err != nil {
			return err
		}
	}
	for _, f := range toRemoveDirs {
		filePath := path.Join(ctx.GetWorkspace().BasePath, mifyconfig.GoServicesRoot, f)
		if err := os.Remove(filePath); err != nil {
			if _, ok := err.(*os.PathError); ok {
				// skip path error because it's likely that dir is not empty
				ctx.Logger.Debugf("not removing %s: %s", filePath, err)
				continue
			}
			return err
		}
	}
	for _, f := range toReplaceLookupPaths {
		filePath := path.Join(ctx.GetWorkspace().BasePath, mifyconfig.GoServicesRoot, f)
		prevSvcPath := fmt.Sprintf("%s/internal/%s/generated", ctx.GetWorkspace().GetGoModule(), ctx.GetServiceName())
		prevPkgPath := fmt.Sprintf("%s/internal/pkg/generated", ctx.GetWorkspace().GetGoModule())
		newSvcPath := fmt.Sprintf("%s/internal/mify-generated/services/%s", ctx.GetWorkspace().GetGoModule(), ctx.GetServiceName())
		newPkgPath := fmt.Sprintf("%s/internal/mify-generated/common", ctx.GetWorkspace().GetGoModule())

		err := filepath.WalkDir(filePath, func(p string, d fs.DirEntry, ferr error) error {
			if d == nil {
				return fmt.Errorf("failed to replace in path %s: %w", p, ferr)
			}
			if d.IsDir() {
				return nil
			}
			if ext := filepath.Ext(p); ext != ".go" {
				return nil
			}
			ctx.Logger.Debugf("replacing %s -> %s in path: %s", prevSvcPath, newSvcPath, p)
			ctx.Logger.Debugf("replacing %s -> %s in path: %s", prevPkgPath, newPkgPath, p)
			migrator.MigrateSubstring(ctx, p, prevSvcPath, "", newSvcPath)
			migrator.MigrateSubstring(ctx, p, prevPkgPath, "", newPkgPath)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateJs(ctx *gencontext.GenContext) error {
	return nil
}

func migratePython(ctx *gencontext.GenContext) error {
	return nil
}
