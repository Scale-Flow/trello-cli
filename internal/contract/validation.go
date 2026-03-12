package contract

import (
	"fmt"
	"sort"
	"strings"
)

// RequireFlag returns a VALIDATION_ERROR if value is empty.
func RequireFlag(name, value string) error {
	if value == "" {
		return NewError(ValidationError, fmt.Sprintf("--%s is required", name))
	}
	return nil
}

// RequireExactlyOne returns a VALIDATION_ERROR if not exactly one flag has a non-empty value.
func RequireExactlyOne(flags map[string]string) error {
	var set []string
	var all []string
	for name, value := range flags {
		all = append(all, "--"+name)
		if value != "" {
			set = append(set, "--"+name)
		}
	}
	sort.Strings(all)
	if len(set) == 1 {
		return nil
	}
	return NewError(ValidationError, fmt.Sprintf("exactly one of %s is required", strings.Join(all, ", ")))
}

// RequireAtLeastOne returns a VALIDATION_ERROR if no flags have a non-empty value.
func RequireAtLeastOne(flags map[string]string) error {
	var all []string
	for name, value := range flags {
		all = append(all, "--"+name)
		if value != "" {
			return nil
		}
	}
	sort.Strings(all)
	return NewError(ValidationError, fmt.Sprintf("at least one of %s is required", strings.Join(all, ", ")))
}
