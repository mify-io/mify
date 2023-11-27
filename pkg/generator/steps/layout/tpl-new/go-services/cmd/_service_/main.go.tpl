{{- .TplHeader }}
// vim: set ft=go:

package main

import (
	"context"
	"{{ .Workspace.MifyGeneratedServicePackage }}/app"
)

func main() {
	app := app.NewMifyServiceApp(context.Background())
	app.Run()
}
