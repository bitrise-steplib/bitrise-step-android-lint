package gradle

import "testing"

func Test_extractArtifactName(t *testing.T) {
	type args struct {
		project Project
		path    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"empty_projpath",
			args{
				project: Project{location: "", monoRepo: true},
				path:    "myartifact.html",
			},
			"myartifact.html",
			false,
		},
		{
			"in_module_projpath",
			args{
				project: Project{location: "root_dir", monoRepo: false},
				path:    "root_dir/mymodule/build/reports/myartifact.html",
			},
			"root_dir-mymodule-myartifact.html",
			false,
		},
		{
			"outside",
			args{
				project: Project{location: "root_dir", monoRepo: true},
				path:    "root_dir/../randomdir/build/reports/myartifact.html",
			},
			"myartifact.html",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractArtifactName(tt.args.project, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractArtifactName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractArtifactName() = %v, want %v", got, tt.want)
			}
		})
	}
}
