package main

import (
	"strings"

	"github.com/bitrise-io/go-utils/log"
)


const resultArtifactPathPattern = "*TEST*.xml"

// getUniqueDir returns the unique subdirectory inside the test addon export directory for a given artifact.
// this assumes the following path is provided: /bitrise/src/[module name]/build/outputs/androidTest-results/connected/flavors/[flavor name]/TEST-xxx.xml
func getUniqueDir(path string) (string, error) {
	log.Debugf("getUniqueDir(%s)", path)
	parts := strings.Split(path, "/")
	
	flavor := parts[len(parts) - 2]

	module := parts[3]
	ret := module + "_" + flavor

	log.Debugf("getUniqueDir(%s): (%s)", path, ret)
	return ret, nil
}
