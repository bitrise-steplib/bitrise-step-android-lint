package gradle

import (
	"reflect"
	"testing"
)

func TestCleanStringSlice(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", []string{"", ""}, []string{}},
		{"one_empty", []string{"", "test"}, []string{"test"}},
		{"one_empty_w/_space", []string{"", "   space "}, []string{"space"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanStringSlice(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
