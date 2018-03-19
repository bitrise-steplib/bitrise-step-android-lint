package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/depman/pathutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/bitrise-steplib/bitrise-step-android-lint/gradle"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	ProjectLocation   string `env:"project_location,dir"`
	ReportPathPattern string `env:"report_path_pattern"`
	Variant           string `env:"variant"`
	Module            string `env:"module"`
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
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
		SetModule(config.Module).
		SetTask("lint")

	log.Infof("Variants:")
	fmt.Println()

	variants, err := lintTask.GetVariants()
	if err != nil {
		failf("Failed to fetch variants, error: %s", err)
	}

	filteredVariants := variants.Filter(config.Variant)

	for _, variant := range variants {
		if sliceutil.IsStringInSlice(variant, filteredVariants) {
			log.Donef("✓ %s", variant)
		} else {
			log.Printf("- %s", variant)
		}
	}

	fmt.Println()

	if len(filteredVariants) == 0 {
		errMsg := fmt.Sprintf("No variant matching for: (%s)", config.Variant)
		if config.Module != "" {
			errMsg += fmt.Sprintf(" in module: [%s]", config.Module)
		}
		failf(errMsg)
	}

	if config.Variant == "" {
		log.Warnf("No variant specified, lint will run on all variants")
		fmt.Println()
	}

	started := time.Now()

	log.Infof("Run lint:")
	if err := lintTask.RunVariants(filteredVariants); err != nil {
		failf("Lint task failed, error: %v", err)
	}
	fmt.Println()

	log.Infof("Exporting artifacts:")
	fmt.Println()

	artifacts, err := gradleProject.FindArtifacts(started, config.ReportPathPattern)
	if err != nil {
		failf("failed to find artifacts, error: %v", err)
	}

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
}
