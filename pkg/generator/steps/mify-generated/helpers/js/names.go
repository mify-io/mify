package js

import (
	"strings"
	"unicode"
)

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
