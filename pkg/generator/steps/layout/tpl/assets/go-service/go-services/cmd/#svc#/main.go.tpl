{{- .Workspace.TplHeader}}
// vim: set ft=go:

package main

import (
	"context"
	"{{ .Workspace.GetAppIncludePath .ServiceName -}}"
)

func main() {
	app := app.NewMifyServiceApp(context.Background())
	app.Run()
}
