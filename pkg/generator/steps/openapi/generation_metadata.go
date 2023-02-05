package openapi

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mify-io/mify/internal/mify/util"
	gencontext "github.com/mify-io/mify/pkg/generator/gen-context"
)

const (
	GENERATION_META_FILENAME = "generation_metadata.yaml"
)

type GenerationMeta struct {
	MifyVersion string `yaml:"mify_version"`
}

func getGenerationMeta(ctx *gencontext.GenContext) (GenerationMeta, error) {
	tmpDir := ctx.GetWorkspace().GetServiceCacheDirectory(ctx.GetServiceName())
	metaFilename := filepath.Join(tmpDir, GENERATION_META_FILENAME)

	var metadata GenerationMeta
	yd := util.NewYAMLData(metaFilename)
	err := yd.ReadFile(&metadata)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return GenerationMeta{}, nil
	}
	if err != nil {
		return GenerationMeta{}, fmt.Errorf("failed to read service generation meta: %w", err)
	}
	return metadata, nil
}

func writeGenerationMeta(ctx *gencontext.GenContext, metadata GenerationMeta) error {
	tmpDir := ctx.GetWorkspace().GetServiceCacheDirectory(ctx.GetServiceName())
	metaFilename := filepath.Join(tmpDir, GENERATION_META_FILENAME)

	yd := util.NewYAMLData(metaFilename)
	err := yd.SaveFile(&metadata)
	if err != nil {
		return fmt.Errorf("failed to save service generation metadata: %w", err)
	}
	return nil

}
