package build

import "runtime"

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func GoVersion() string { return runtime.Version() }
