package gradle

import "testing"

func TestGetGradleModule(t *testing.T) {
	tests := []struct {
		name         string
		configModule string
		want         string
	}{
		{"empty", "", ""},
		{"simple", "app", ":app:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGradleModule(tt.configModule); got != tt.want {
				t.Errorf("getGradleModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
