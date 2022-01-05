{{- .Workspace.TplHeader}}

package main

import (
	"context"
	"{{ .Workspace.GetAppIncludePath .ServiceName -}}"
)

func main() {
	app := app.NewMifyServiceApp(context.Background())
	app.Run()
}
