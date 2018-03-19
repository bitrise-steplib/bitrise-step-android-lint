package gradle

import (
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
)

// Artifact ...
type Artifact struct {
	Path string
	Name string
}

func extractArtifactName(project Project, path string) (string, error) {
	location := project.location

	if project.monoRepo {
		location = filepath.Join(project.location, "..")
	}

	relPath, err := filepath.Rel(location, path)
	if err != nil {
		return "", err
	}

	s := strings.Split(relPath, "/")
	suffix := s[len(s)-1]

	// return with filename only if user decided to have
	// the result files outside of the project dir
	if !strings.HasPrefix(relPath, "..") && project.location != "." && project.location != "" {
		s := strings.Split(relPath, "/")
		module := s[0]

		if project.monoRepo && len(s) > 1 {
			module += "-" + s[1]
		}

		return module + "-" + suffix, nil
	}

	return suffix, nil
}

// Export ...
func (artifact Artifact) Export(destination string) error {
	return command.CopyFile(artifact.Path, filepath.Join(destination, artifact.Name))
}
