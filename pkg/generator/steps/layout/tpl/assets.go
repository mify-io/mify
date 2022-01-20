package tpl

import "embed"

//go:embed *
var tpl embed.FS

func GetTplFs() embed.FS {
	return tpl
}
