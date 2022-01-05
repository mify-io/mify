{{- .Workspace.TplHeader}}

package main

import (
	"context"
{{- range $key, $value := .Workspace.GoServices }}
	{{ $value.Name }} "{{ $.Workspace.GetAppIncludePath $value.Name -}}"
{{- end }}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
{{ range $key, $value := .Workspace.GoServices }}
	app_{{ $value.Name }} := {{ $value.Name }}.NewMifyServiceApp(ctx)
	go func() {
		app_{{ $value.Name }}.Run()
		cancel()
	}()
{{ end }}
	<-ctx.Done()
}
