package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bitrise-io/go-android/gradle"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"github.com/zawadz88/bitrise-step-junit-results-exporter/testaddon"
)

// Configs ...
type Configs struct {
	ProjectLocation   			string `env:"project_location,dir"`
	TestFolderName 				string `env:"test_folder_name"`
	ResultArtifactPathPattern 	string `env:"result_artifact_path_pattern"`
}

func failf(f string, args ...interface{}) {
	log.Errorf(f, args...)
	os.Exit(1)
}

func getArtifacts(gradleProject gradle.Project, started time.Time, pattern string) (artifacts []gradle.Artifact, err error) {
	for _, t := range []time.Time{started, time.Time{}} {
		artifacts, err = gradleProject.FindArtifacts(t, pattern, false)
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
	var config Configs

	if err := stepconf.Parse(&config); err != nil {
		failf("Couldn't create step config: %v\n", err)
	}

	stepconf.Print(config)
	fmt.Println()

	gradleProject, err := gradle.NewProject(config.ProjectLocation)
	if err != nil {
		failf("Failed to open project, error: %s", err)
	}

	started := time.Now().Add(time.Duration(-90) * time.Minute)

	log.Infof("Export test results for test addon:")
	fmt.Println()

	resultXMLs, err := getArtifacts(gradleProject, started, config.ResultArtifactPathPattern)
	if err != nil {
		log.Warnf("Failed to find test result XMLs, error: %s", err)
	} else {
		if baseDir := os.Getenv("BITRISE_TEST_RESULT_DIR"); baseDir != "" {
			for _, artifact := range resultXMLs {
				uniqueDir, err := getUniqueDir(artifact.Path, config.TestFolderName)
				if err != nil {
					log.Warnf("Failed to export test results for test addon: cannot get export directory for artifact (%s): %s", artifact.Name, err)
					continue
				}
				log.Printf("  Exporting artifact to test addon[ path: %s, name: %s, uniqueDir: %s ]", artifact.Path, artifact.Name, uniqueDir)

				if err := testaddon.ExportArtifact(artifact.Path, baseDir, uniqueDir); err != nil {
					log.Warnf("Failed to export test results for test addon: %s", err)
				}
			}
			log.Printf("  Exporting test results to test addon successful [ %s ] ", baseDir)
		}
	}

	fmt.Println()
	log.Donef("  Done")
}
