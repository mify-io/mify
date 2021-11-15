package config

import "path/filepath"

// import (
// "github.com/adrg/xdg"
// )
const (
	TMP_SUBDIR          = "services"
)

func GetCacheDirectory(basePath string) string {
	// return xdg.CacheHome+"/mify"
	return filepath.Join(basePath, ".mify")
}

func GetServiceCacheDirectory(basePath, serviceName string) string {
	return filepath.Join(GetCacheDirectory(basePath), TMP_SUBDIR, serviceName)
}
