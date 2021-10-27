package config

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/chebykinn/mify/assets"
)

func GetAssets() embed.FS {
	return assets.GetAssetsFs()
}

func DumpAssets(assetPath string, targetDir string) (string, error) {
	assetsFs := GetAssets()
	cacheDir := GetCacheDirectory()
	targetDir = filepath.Join(cacheDir, "assets", targetDir)

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return "", err
	}

	err = copyImpl(assetsFs, assetPath, targetDir)
	if err != nil {
		return "", err
	}

	return targetDir, nil
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
