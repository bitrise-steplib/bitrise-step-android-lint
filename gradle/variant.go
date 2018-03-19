package gradle

import (
	"strings"

	"github.com/bitrise-io/go-utils/sliceutil"
)

// Variants ...
type Variants []string

// Filter ...
func (variants Variants) Filter(filter string) Variants {
	cleanedFilters := cleanStringSlice(strings.Split(filter, "\n"))

	if len(cleanedFilters) == 0 {
		return variants
	}

	t := []string{}
	for _, v := range variants {
		for _, f := range cleanedFilters {
			f = strings.ToLower(f)
			if strings.Contains(strings.ToLower(v), f) {
				if !sliceutil.IsStringInSlice(v, t) {
					t = append(t, v)
				}
			}
		}
	}
	return t
}
