package main

import "ai-text-shaper/cmd"

var (
	// Version is the version of the application.
	version = "dev"
	// Commit is the git commit of the build.
	commit = "none"
	// Date is the date of the build.
	date = "unknown"
	// BuiltBy is the user/tool that built the binary.
	builtBy = "dirty hands"
)

func main() {
	cmd.Execute(version, commit, date, builtBy)
}
