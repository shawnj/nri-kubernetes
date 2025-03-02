package testutil

import (
	"embed"
)

//go:embed data
var testDataDir embed.FS

// Name of the root folder in embed.FS
const testDataRootDir = "data"

// Version represents a kubernetes version. Mock servers can be instantiated to return known output for a given version.
type Version string

// Server returns an HTTP Server for the given version, ready to serve static endpoints for KSM, Kubelet and CP components.
func (v Version) Server() (*Server, error) {
	return newServer(v)
}

// K8s returns a helper that provide fake instances of K8s objects, ready to use with the kubernetes fake client.
func (v Version) K8s() (K8s, error) {
	return newK8s(v)
}

// List of all the versions we have testdata for.
// When adding a new version:
// - REMEMBER TO ADD IT TO AllVersions() BELOW.

const (
	Testdata116 = "1_16"
	Testdata117 = "1_17"
	Testdata118 = "1_18"
	Testdata119 = "1_19"
	Testdata120 = "1_20"
	Testdata121 = "1_21"
	Testdata122 = "1_22"
)

// AllVersions returns a list of versions we have test data for.
// PLEASE ADD NEW VERSIONS HERE AS WELL.
// PLEASE KEEP THIS LIST SORTED, WITH NEWER RELEASES LAST IN THE LIST.
func AllVersions() []Version {
	return []Version{
		Testdata116,
		Testdata117,
		Testdata118,
		Testdata119,
		Testdata120,
		Testdata121,
		Testdata122,
	}
}

// LatestVersion returns the latest version we have test data for.
func LatestVersion() Version {
	allVersions := AllVersions()
	return allVersions[len(allVersions)-1]
}
