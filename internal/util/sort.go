package util

import (
	"sort"

	"github.com/fsrv-xyz/nyx/internal/check"
)

func SortByCheckName(checks []check.GenericCheck) []check.GenericCheck {
	mapping := make(map[string]check.GenericCheck)
	keys := make([]string, 0, len(checks))
	result := make([]check.GenericCheck, 0, len(checks))

	for checkIndex := range checks {
		mapping[checks[checkIndex].Name] = checks[checkIndex]
		keys = append(keys, checks[checkIndex].Name)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result = append(result, mapping[k])
	}
	return result
}
