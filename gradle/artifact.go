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
	relPath, err := filepath.Rel(project.location, path)
	if err != nil {
		return "", err
	}

	module := strings.Split(relPath, "/")[0]
	fileName := filepath.Base(relPath)

	if project.monoRepo {
		splitPath := strings.Split(project.location, "/")
		module = splitPath[len(splitPath)-1] + "-" + module
	}

	return module + "-" + fileName, nil
}

// Export ...
func (artifact Artifact) Export(destination string) error {
	return command.CopyFile(artifact.Path, filepath.Join(destination, artifact.Name))
}
