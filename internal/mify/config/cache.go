package config

import "path/filepath"

// import (
// "github.com/adrg/xdg"
// )

func GetCacheDirectory(basePath string) string {
	// return xdg.CacheHome+"/mify"
	return filepath.Join(basePath, ".mify")
}
