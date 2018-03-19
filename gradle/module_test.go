package gradle

import "testing"

func Test_getGradleModule(t *testing.T) {
	type args struct {
		configModule string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{configModule: ""}, ""},
		{"simple", args{configModule: "app"}, ":app:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGradleModule(tt.args.configModule); got != tt.want {
				t.Errorf("getGradleModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
