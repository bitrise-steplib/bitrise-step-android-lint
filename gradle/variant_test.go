package gradle

import (
	"reflect"
	"testing"
)

func TestVariantsFilter(t *testing.T) {
	tests := []struct {
		name     string
		variants Variants
		filter   string
		want     Variants
	}{
		{
			"empty",
			sampleVariants,
			``,
			sampleVariants,
		},
		{
			"only_newlines",
			sampleVariants,
			`
				
				
				
				`,
			sampleVariants,
		},
		{
			"only_match",
			sampleVariants,
			`stage`,
			Variants{
				"InvArm7StageDebug",
				"InvArm7StageRelease",
				"InvX86StageDebug",
				"InvX86StageRelease",
			},
		},
		{
			"match_and_newlines",
			sampleVariants,
			`
				stage
				
				
				`,
			Variants{
				"InvArm7StageDebug",
				"InvArm7StageRelease",
				"InvX86StageDebug",
				"InvX86StageRelease",
			},
		},
		{
			"match_multiple_times_and_newlines",
			sampleVariants,
			`
				stage
				
				stage`,
			Variants{
				"InvArm7StageDebug",
				"InvArm7StageRelease",
				"InvX86StageDebug",
				"InvX86StageRelease",
			},
		},
		{
			"multiple_match_multiple_times_and_newlines",
			sampleVariants,
			`
				stage
				
				stage
				staging`,
			Variants{
				"Myflavor2Staging",
				"MyflavorokStaging",
				"MyflavorStaging",
				"InvArm7StageDebug",
				"InvArm7StageRelease",
				"InvX86StageDebug",
				"InvX86StageRelease",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.variants.Filter(tt.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variants.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

var sampleVariants = Variants{
	"Myflavor2Debug",
	"Myflavor2Release",
	"Myflavor2Staging",
	"MyflavorDebug",
	"MyflavorokDebug",
	"MyflavorokRelease",
	"MyflavorokStaging",
	"MyflavorRelease",
	"MyflavorStaging",
	"InvArm7LocalDebug",
	"InvArm7LocalRelease",
	"InvArm7ProdDebug",
	"InvArm7ProdRelease",
	"InvArm7StageDebug",
	"InvArm7StageRelease",
	"InvX86LocalDebug",
	"InvX86LocalRelease",
	"InvX86ProdDebug",
	"InvX86ProdRelease",
	"InvX86StageDebug",
	"InvX86StageRelease",
}
