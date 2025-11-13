package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitrise-io/go-android/cache"
	"github.com/bitrise-io/go-android/gradle"
	utilscache "github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/env"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"github.com/kballard/go-shellquote"
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

var logger = log.NewLogger(false)

func getArtifacts(gradleProject gradle.Project, started time.Time, pattern string) (artifacts []gradle.Artifact, err error) {
	artifacts, err = gradleProject.FindArtifacts(started, pattern, true)
	if err != nil {
		return
	}
	if len(artifacts) == 0 {
		if !started.IsZero() {
			logger.Warnf("No artifacts found with pattern: %s that has modification time after: %s", pattern, started)
			logger.Warnf("Retrying without modtime check....")
			logger.Println()
			return getArtifacts(gradleProject, time.Time{}, pattern)
		}
		logger.Warnf("No artifacts found with pattern: %s without modtime check", pattern)
		logger.Warnf("If you have changed default report export path in your gradle files then you might need to change ReportPathPattern accordingly.")
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
			if strings.EqualFold(v, variant) {
				filteredVariants[m] = append(filteredVariants[m], v)
			}
		}
	}
	if len(filteredVariants) == 0 {
		return nil, fmt.Errorf("variant: %s not found in any module", variant)
	}
	return filteredVariants, nil
}

func mainE(config Config, cmdFactory command.Factory, logger log.Logger) error {
	gradleProject, err := gradle.NewProject(config.ProjectLocation, cmdFactory)
	if err != nil {
		return fmt.Errorf("Process config: failed to open project, error: %s", err)
	}

	lintTask := gradleProject.GetTask("lint")

	args, err := shellquote.Split(config.Arguments)
	if err != nil {
		return fmt.Errorf("Process config: failed to parse arguments, error: %s", err)
	}

	logger.Infof("Variants:")
	fmt.Println()

	variants, err := lintTask.GetVariants(args...)
	if err != nil {
		return fmt.Errorf("Run: failed to fetch variants, error: %s", err)
	}

	filteredVariants, err := filterVariants(config.Module, config.Variant, variants)
	if err != nil {
		failf("Process config: failed to find buildable variants, error: %s", err)
	}

	for module, variants := range variants {
		logger.Printf("%s:", module)
		for _, variant := range variants {
			if sliceutil.IsStringInSlice(variant, filteredVariants[module]) {
				logger.Donef("✓ %s", strings.TrimSuffix(variant, "UnitTest"))
				continue
			}
			logger.Printf("- %s", strings.TrimSuffix(variant, "UnitTest"))
		}
	}
	fmt.Println()

	started := time.Now()

	logger.Infof("Run lint:")
	lintCommand := lintTask.GetCommand(filteredVariants, args...)

	fmt.Println()
	logger.Donef("$ " + lintCommand.PrintableCommandArgs())
	fmt.Println()

	taskError := lintCommand.Run()
	if taskError != nil {
		logger.Errorf("Run: lint task failed, error: %v", taskError)
	}
	fmt.Println()

	logger.Infof("Exporting artifacts:")
	fmt.Println()

	artifacts, err := getArtifacts(gradleProject, started, config.ReportPathPattern)
	if err != nil {
		return fmt.Errorf("Export outputs: failed to find artifacts, error: %v", err)
	}

	if len(artifacts) > 0 {
		for _, artifact := range artifacts {
			exists, err := pathutil.IsPathExists(
				filepath.Join(config.DeployDir, artifact.Name),
			)
			if err != nil {
				return fmt.Errorf("Export outputs: failed to check path, error: %v", err)
			}

			artifactName := filepath.Base(artifact.Path)

			if exists {
				timestamp := time.Now().
					Format("20060102150405")
				ext := filepath.Ext(artifact.Name)
				name := strings.TrimSuffix(filepath.Base(artifact.Name), ext)
				artifact.Name = fmt.Sprintf("%s-%s%s", name, timestamp, ext)
			}

			logger.Printf("  Export [ %s => $BITRISE_DEPLOY_DIR/%s ]", artifactName, artifact.Name)

			if err := artifact.Export(config.DeployDir); err != nil {
				logger.Warnf("failed to export artifacts, error: %v", err)
			}
		}
	} else {
		logger.Warnf("No artifacts found with pattern: %s", config.ReportPathPattern)
		logger.Warnf("If you have changed default report file paths with lintOptions/htmlOutput or lintOptions/xmlOutput")
		logger.Warnf("in your gradle files then you might need to change ReportPathPattern accordingly.")
	}

	return taskError
}

func failf(f string, args ...interface{}) {
	logger.Errorf(f, args...)
	os.Exit(1)
}

func main() {
	var config Config

	if err := stepconf.Parse(&config); err != nil {
		failf("Process config: couldn't create step config: %v\n", err)
	}

	stepconf.Print(config)
	fmt.Println()

	cmdFactory := command.NewFactory(env.NewRepository())

	if err := mainE(config, cmdFactory, logger); err != nil {
		failf("%s", err)
	}

	fmt.Println()
	logger.Infof("Collecting cache:")
	if warning := cache.Collect(config.ProjectLocation, utilscache.Level(config.CacheLevel), cmdFactory); warning != nil {
		logger.Warnf("%s", warning)
	}
	logger.Donef("  Done")
}
