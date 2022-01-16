package service

import (
	"strings"
	"unicode"

	"github.com/chebykinn/mify/internal/mify/workspace"
	"github.com/chebykinn/mify/pkg/mifyconfig"
)

type Context struct {
	ServiceName string
	Repository  string
	Language    mifyconfig.ServiceLanguage
	GoModule    string
	Workspace   workspace.Context
	ServiceList []string
}

func (c Context) GetEndpointEnvName() string {
	return MakeServerEnvName(c.ServiceName)
}

func SanitizeServiceName(serviceName string) string {
	if unicode.IsDigit(rune(serviceName[0])) {
		serviceName = "service_" + serviceName
	}
	serviceName = strings.ReplaceAll(serviceName, "-", "_")

	return serviceName
}

func MakeServerEnvName(serviceName string) string {
	sanitizedName := SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_SERVER_ENDPOINT"
}
