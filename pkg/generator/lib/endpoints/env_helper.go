package endpoints

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

func MakeApiEndpointEnvName(serviceName string) string {
	sanitizedName := SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_API_ENDPOINT"
}

func MakeMaintenanceEndpointEnvName(serviceName string) string {
	sanitizedName := SanitizeServiceName(serviceName)
	return strings.ToUpper(sanitizedName) + "_MAINTENANCE_ENDPOINT"
}

func SnakeCaseToCamelCase(inputUnderScoreStr string, capitalize bool) (camelCase string) {
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 && capitalize {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}
