package constant

import "runtime"

// Owner, Repo, OS and Arch identify the source repository and host platform,
// used when checking for and downloading self-update releases.
var (
	Owner = "sjf10050"
	Repo  = "nali"
	OS    = runtime.GOOS
	Arch  = runtime.GOARCH
)
