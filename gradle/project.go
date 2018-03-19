package gradle

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/bitrise-io/depman/pathutil"
	"github.com/bitrise-io/go-utils/log"
	glob "github.com/ryanuber/go-glob"
)

// Project ...
type Project struct {
	location string
	monoRepo bool
}

// NewProject ...
func NewProject(location string) (Project, error) {
	buildGradleFound, err := pathutil.IsPathExists(filepath.Join(location, "build.gradle"))
	if err != nil {
		return Project{}, err
	}

	if !buildGradleFound {
		return Project{}, fmt.Errorf("no build.gradle file found in (%s)", location)
	}

	root := filepath.Join(location, "..")

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return Project{}, err
	}

	projectsCount := 0
	for _, file := range files {
		if file.IsDir() {
			e, err := pathutil.IsPathExists(filepath.Join(root, file.Name(), "build.gradle"))
			if err != nil {
				return Project{}, err
			}
			if e {
				projectsCount++
			}
		}
	}

	return Project{location: location, monoRepo: (projectsCount < 2)}, nil
}

// SetModule ...
func (proj Project) SetModule(module string) Module {
	return Module{
		project: proj,
		name:    getGradleModule(module),
	}
}

// FindArtifacts ...
func (proj Project) FindArtifacts(generatedAfter time.Time, pattern string) ([]Artifact, error) {
	paths := []Artifact{}
	return paths, filepath.Walk(proj.location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("failed to walk path: %s", err)
			return nil
		}

		if info.ModTime().Before(generatedAfter) || !glob.Glob(pattern, path) || info.IsDir() {
			return nil
		}

		name, err := extractArtifactName(proj, path)
		if err != nil {
			return err
		}

		paths = append(paths, Artifact{Name: name, Path: path})
		return nil
	})
}
