package gradle

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
)

// If we parse tasks that starts with lint, we will have tasks that starts
// with lintVital also. So list here each conflicting tasks. (only overlapping ones)
var conflicts = map[string][]string{
	"lint": []string{
		"lintVital",
	},
}

func getGradleOutput(projPath string, tasks ...string) (string, error) {
	c := command.New(filepath.Join(projPath, "gradlew"), tasks...)
	return c.RunAndReturnTrimmedCombinedOutput()
}

func runGradleCommand(projPath string, tasks ...string) error {
	return command.New(filepath.Join(projPath, "gradlew"), tasks...).
		SetDir(projPath).
		SetStdout(os.Stdout).
		SetStderr(os.Stderr).
		Run()
}

func cleanStringSlice(in []string) (out []string) {
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return
}
