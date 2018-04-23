package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-steplib/bitrise-step-android-unit-test/cache"
	"github.com/bitrise-tools/go-android/gradle"
	"github.com/bitrise-tools/go-steputils/stepconf"
	shellquote "github.com/kballard/go-shellquote"
)

// Config ...
type Config struct {
	ProjectLocation   string `env:"project_location,dir"`
	ReportPathPattern string `env:"report_path_pattern"`
	Variant           string `env:"variant"`
	Module            string `env:"module"`
	Arguments         string `env:"arguments"`
	CacheLevel        string `env:"cache_level,opt[none,only_deps,all]"`
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

func getArtifacts(gradleProject gradle.Project, started time.Time, pattern string) (artifacts []gradle.Artifact, err error) {
	for _, t := range []time.Time{started, time.Time{}} {
		artifacts, err = gradleProject.FindArtifacts(t, pattern, true)
		if err != nil {
			return
		}
		if len(artifacts) == 0 {
			if t == started {
				log.Warnf("No artifacts found with pattern: %s that has modification time after: %s", pattern, t)
				log.Warnf("Retrying without modtime check....")
				fmt.Println()
				continue
			}
			log.Warnf("No artifacts found with pattern: %s without modtime check", pattern)
			log.Warnf("If you have changed default report export path in your gradle files then you might need to change ReportPathPattern accordingly.")
		}
	}
	return
}

func main() {
	var config Config

	if err := stepconf.Parse(&config); err != nil {
		failf("Couldn't create step config: %v\n", err)
	}

	stepconf.Print(config)

	deployDir := os.Getenv("BITRISE_DEPLOY_DIR")

	log.Printf("- Deploy dir: %s", deployDir)
	fmt.Println()

	gradleProject, err := gradle.NewProject(config.ProjectLocation)
	if err != nil {
		failf("Failed to open project, error: %s", err)
	}

	lintTask := gradleProject.
		GetTask("lint")

	log.Infof("Variants:")
	fmt.Println()

	variants, err := lintTask.GetVariants()
	if err != nil {
		failf("Failed to fetch variants, error: %s", err)
	}

	filteredVariants := variants.Filter(config.Module, config.Variant)

	for module, variants := range variants {
		log.Printf("%s:", module)
		for _, variant := range variants {
			if sliceutil.IsStringInSlice(variant, filteredVariants[module]) {
				log.Donef("✓ %s", strings.TrimSuffix(variant, "UnitTest"))
			} else {
				log.Printf("- %s", strings.TrimSuffix(variant, "UnitTest"))
			}
		}
	}
	fmt.Println()

	if len(filteredVariants) == 0 {
		if config.Variant != "" {
			if config.Module == "" {
				failf("Variant (%s) not found in any module", config.Variant)
			} else {
				failf("No variant matching for (%s) in module: [%s]", config.Variant, config.Module)
			}
		}
		failf("Module not found: %s", config.Module)
	}

	started := time.Now()

	args, err := shellquote.Split(config.Arguments)
	if err != nil {
		failf("Failed to parse arguments, error: %s", err)
	}

	log.Infof("Run lint:")
	taskError := lintTask.Run(filteredVariants, args...)
	if taskError != nil {
		log.Errorf("Lint task failed, error: %v", err)
	}
	fmt.Println()

	log.Infof("Exporting artifacts:")
	fmt.Println()

	artifacts, err := getArtifacts(gradleProject, started, config.ReportPathPattern)
	if err != nil {
		failf("failed to find artifacts, error: %v", err)
	}

	if len(artifacts) > 0 {
		for _, artifact := range artifacts {
			exists, err := pathutil.IsPathExists(
				filepath.Join(deployDir, artifact.Name),
			)
			if err != nil {
				failf("failed to check path, error: %v", err)
			}

			artifactName := filepath.Base(artifact.Path)

			if exists {
				timestamp := time.Now().
					Format("20060102150405")
				ext := filepath.Ext(artifact.Name)
				name := strings.TrimSuffix(filepath.Base(artifact.Name), ext)
				artifact.Name = fmt.Sprintf("%s-%s%s", name, timestamp, ext)
			}

			log.Printf("  Export [ %s => $BITRISE_DEPLOY_DIR/%s ]", artifactName, artifact.Name)

			if err := artifact.Export(deployDir); err != nil {
				log.Warnf("failed to export artifacts, error: %v", err)
			}
		}
	} else {
		log.Warnf("No artifacts found with pattern: %s", config.ReportPathPattern)
		log.Warnf("If you have changed default report file paths with lintOptions/htmlOutput or lintOptions/xmlOutput")
		log.Warnf("in your gradle files then you might need to change ReportPathPattern accordingly.")
	}

	if taskError != nil {
		os.Exit(1)
	}

	fmt.Println()
	log.Infof("Collecting cache:")
	if warning := cache.Collect(config.ProjectLocation, cache.Level(config.CacheLevel)); warning != nil {
		log.Warnf("%s", warning)
	}
	log.Donef("  Done")
}
