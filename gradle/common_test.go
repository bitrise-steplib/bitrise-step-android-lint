package gradle

import (
	"reflect"
	"testing"
)

func Test_cleanStringSlice(t *testing.T) {
	type args struct {
		in []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty", args{in: []string{"", ""}}, []string{}},
		{"one_empty", args{in: []string{"", "test"}}, []string{"test"}},
		{"one_empty_w/_space", args{in: []string{"", "   space "}}, []string{"space"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanStringSlice(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
