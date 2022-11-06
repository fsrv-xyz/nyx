package util

import (
	"sort"

	"github.com/fsrv-xyz/nyx/internal/check"
)

func SortByCheckName(checks []check.GenericCheck) []check.GenericCheck {
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Name < checks[j].Name
	})
	return checks
}
