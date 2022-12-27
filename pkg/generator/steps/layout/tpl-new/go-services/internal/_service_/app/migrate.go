package app

import (
	"strings"
)

// #65
func migrateContextToServiceExtra(_ string, _ interface{}, currentText string) (string, error) {
	// TODO: use https://pkg.go.dev/go/parser
	return strings.ReplaceAll(
		currentText,
		"func NewRouterConfig()",
		"func NewRouterConfig(ctx *core.MifyServiceContext)"), nil
}
