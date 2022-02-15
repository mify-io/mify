package util

import (
	"fmt"
	"strings"
)

func ValidateStrArg(arg string, allowed []string) error {
	found := false
	for _, t := range allowed {
		if arg == t {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf(
			"value '%s' is not allowed. List of possible args: [%s]",
			arg, strings.Join(allowed, ","))
	}

	return nil
}
