{{- .Workspace.TplHeader}}

package main

import (
	"{{.Repository}}/go_services/internal/{{.ServiceName}}/generated/app"
)

func main() {
	app := app.NewMifyServiceApp()
	app.Run()
}
