{{- .TplHeader }}
// vim: set ft=go:

package main

import (
	"context"
	"{{ .Workspace.PackageName }}/internal/{{ .ServiceName }}/generated/app"
)

func main() {
	app := app.NewMifyServiceApp(context.Background())
	app.Run()
}
