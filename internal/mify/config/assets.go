package config

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mify-io/mify/assets"
)

func GetAssets() embed.FS {
	return assets.GetAssetsFs()
}

func HasAssets(assetPath string) bool {
	assetsFs := GetAssets()
	_, err := assetsFs.Open(assetPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func DumpAssets(basePath string, assetPath string, targetDir string) (string, error) {
	assetsFs := GetAssets()
	cacheDir := GetCacheDirectory(basePath)
	targetDir = filepath.Join(cacheDir, "assets", targetDir)

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return "", err
	}

	err = copyImpl(assetsFs, assetPath, targetDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, filepath.Base(assetPath)), nil
}

func copyImpl(fs embed.FS, curPath string, targetDir string) error {
	f, err := fs.Open(curPath)
	if err != nil {
		return err
	}

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		subTargetDir := filepath.Join(targetDir, stat.Name())
		err = os.MkdirAll(subTargetDir, 0755)
		if err != nil {
			return err
		}

		files, err := fs.ReadDir(curPath)
		if err != nil {
			return err
		}
		for _, ent := range files {
			subPath := filepath.Join(curPath, ent.Name())
			if err = copyImpl(fs, subPath, subTargetDir); err != nil {
				return err
			}
		}
		return nil
	}
	data, err := fs.ReadFile(curPath)
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(targetDir, stat.Name()), data, 0644); err != nil {
		return err
	}

	return nil
}
