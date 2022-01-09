package util

import "strings"

func ToSafeGoVariableName(name string) string {
	return strings.Replace(name, "-", "_", -1)
}
