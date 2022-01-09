{{- .Workspace.TplHeader}}

package main

import (
	"context"
{{- range $key, $value := .Workspace.GoServices }}
	{{ $value.GetSafeName }} "{{ $.Workspace.GetAppIncludePath $value.Name -}}"
{{- end }}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
{{ range $key, $value := .Workspace.GoServices }}
	app_{{ $value.GetSafeName }} := {{ $value.GetSafeName }}.NewMifyServiceApp(ctx)
	go func() {
		app_{{ $value.GetSafeName }}.Run()
		cancel()
	}()
{{ end }}
	<-ctx.Done()
}
