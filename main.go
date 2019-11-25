package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-android/cache"
	"github.com/bitrise-io/go-android/gradle"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
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
	DeployDir         string `env:"BITRISE_DEPLOY_DIR,dir"`
}

func getArtifacts(gradleProject gradle.Project, started time.Time, pattern string) (artifacts []gradle.Artifact, err error) {
	artifacts, err = gradleProject.FindArtifacts(started, pattern, true)
	if err != nil {
		return
	}
	if len(artifacts) == 0 {
		if !started.IsZero() {
			log.Warnf("No artifacts found with pattern: %s that has modification time after: %s", pattern, started)
			log.Warnf("Retrying without modtime check....")
			fmt.Println()
			return getArtifacts(gradleProject, time.Time{}, pattern)
		}
		log.Warnf("No artifacts found with pattern: %s without modtime check", pattern)
		log.Warnf("If you have changed default report export path in your gradle files then you might need to change ReportPathPattern accordingly.")
	}
	return
}

func filterVariants(module, variant string, variantsMap gradle.Variants) (gradle.Variants, error) {
	// if module set: drop all the other modules
	if module != "" {
		v, ok := variantsMap[module]
		if !ok {
			return nil, fmt.Errorf("module not found: %s", module)
		}
		variantsMap = gradle.Variants{module: v}
	}
	// if variant not set: use all variants
	if variant == "" {
		return variantsMap, nil
	}
	filteredVariants := gradle.Variants{}
	for m, variants := range variantsMap {
		for _, v := range variants {
			if strings.ToLower(v) == strings.ToLower(variant) {
				filteredVariants[m] = append(filteredVariants[m], v)
			}
		}
	}
	if len(filteredVariants) == 0 {
		return nil, fmt.Errorf("variant: %s not found in any module", variant)
	}
	return filteredVariants, nil
}

func mainE(config Config) error {
	gradleProject, err := gradle.NewProject(config.ProjectLocation)
	if err != nil {
		return fmt.Errorf("Failed to open project, error: %s", err)
	}

	lintTask := gradleProject.
		GetTask("lint")

	log.Infof("Variants:")
	fmt.Println()

	variants, err := lintTask.GetVariants()
	if err != nil {
		return fmt.Errorf("Failed to fetch variants, error: %s", err)
	}

	filteredVariants, err := filterVariants(config.Module, config.Variant, variants)
	if err != nil {
		failf("Failed to find buildable variants, error: %s", err)
	}

	for module, variants := range variants {
		log.Printf("%s:", module)
		for _, variant := range variants {
			if sliceutil.IsStringInSlice(variant, filteredVariants[module]) {
				log.Donef("✓ %s", strings.TrimSuffix(variant, "UnitTest"))
				continue
			}
			log.Printf("- %s", strings.TrimSuffix(variant, "UnitTest"))
		}
	}
	fmt.Println()

	started := time.Now()

	args, err := shellquote.Split(config.Arguments)
	if err != nil {
		return fmt.Errorf("Failed to parse arguments, error: %s", err)
	}

	log.Infof("Run lint:")
	lintCommand := lintTask.GetCommand(filteredVariants, args...)

	fmt.Println()
	log.Donef("$ " + lintCommand.PrintableCommandArgs())
	fmt.Println()

	taskError := lintCommand.Run()
	if taskError != nil {
		log.Errorf("Lint task failed, error: %v", taskError)
	}
	fmt.Println()

	log.Infof("Exporting artifacts:")
	fmt.Println()

	artifacts, err := getArtifacts(gradleProject, started, config.ReportPathPattern)
	if err != nil {
		return fmt.Errorf("failed to find artifacts, error: %v", err)
	}

	if len(artifacts) > 0 {
		for _, artifact := range artifacts {
			exists, err := pathutil.IsPathExists(
				filepath.Join(config.DeployDir, artifact.Name),
			)
			if err != nil {
				return fmt.Errorf("failed to check path, error: %v", err)
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

			if err := artifact.Export(config.DeployDir); err != nil {
				log.Warnf("failed to export artifacts, error: %v", err)
			}
		}
	} else {
		log.Warnf("No artifacts found with pattern: %s", config.ReportPathPattern)
		log.Warnf("If you have changed default report file paths with lintOptions/htmlOutput or lintOptions/xmlOutput")
		log.Warnf("in your gradle files then you might need to change ReportPathPattern accordingly.")
	}

	return taskError
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
	fmt.Println()

	if err := mainE(config); err != nil {
		failf("%s", err)
	}

	fmt.Println()
	log.Infof("Collecting cache:")
	if warning := cache.Collect(config.ProjectLocation, cache.Level(config.CacheLevel)); warning != nil {
		log.Warnf("%s", warning)
	}
	log.Donef("  Done")
}
