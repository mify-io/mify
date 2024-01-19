{{- .TplHeader}}
// vim: set ft=go:

package apputil

import (
	"{{.Model.AppImportPath}}"
	"{{.Model.CoreImportPath}}"
)

func GetServiceExtra(ctx *core.MifyServiceContext) *app.ServiceExtra {
	return ctx.ServiceExtra().(*app.ServiceExtra)
}

func GetRequestExtra(ctx *core.MifyRequestContext) *app.RequestExtra {
	return ctx.RequestExtra().(*app.RequestExtra)
}
