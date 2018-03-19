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

	cleanedVariants := []string{}
	for _, variant := range variants {
		for _, filter := range cleanedFilters {
			filter = strings.ToLower(filter)
			if strings.Contains(strings.ToLower(variant), filter) {
				if !sliceutil.IsStringInSlice(variant, cleanedVariants) {
					cleanedVariants = append(cleanedVariants, variant)
				}
			}
		}
	}
	return cleanedVariants
}
