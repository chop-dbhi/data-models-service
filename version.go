package main

import (
	"strconv"
	"time"

	"github.com/blang/semver"
)

// SemVer components.
const (
	progMajor      = 1
	progMinor      = 0
	progPatch      = 0
	progRelease    = "beta"
	progReleaseNum = 1
)

var (
	// Populated at build time. See the Makefile for details.
	progBuild     string
	progTimestamp string

	// Parsed date from timestamp.
	progDate time.Time

	// Full semantic version for the service.
	progVersion = semver.Version{
		Major: progMajor,
		Minor: progMinor,
		Patch: progPatch,
		Pre: []semver.PRVersion{{
			VersionStr: progRelease,
			VersionNum: progReleaseNum,
		}},
		Build: []string{progBuild},
	}
)

func init() {
	// Convert build timestamp to date.
	ts, _ := strconv.ParseInt(progTimestamp, 10, 64)
	progDate = time.Unix(ts, 0).UTC()
}
