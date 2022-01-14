{{- .Header}}

package main

{{if .Services }}

import (
	"context"
{{- range $key, $value := .Services }}
	{{ $value.SafeName }} "{{ $value.AppIncludePath -}}"
{{- end }}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
{{ range $key, $value := .Services }}
	app_{{ $value.SafeName }} := {{ $value.SafeName }}.NewMifyServiceApp(ctx)
	go func() {
		app_{{ $value.SafeName }}.Run()
		cancel()
	}()
{{ end }}
	<-ctx.Done()
}

{{else}}
func main() {
}
{{end}}
