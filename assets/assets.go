package assets

import "embed"

//go:embed *
var assets embed.FS

func GetAssetsFs() embed.FS {
	return assets
}
