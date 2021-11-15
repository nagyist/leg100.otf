package otf

import "strconv"

var (
	// Build-time parameters set -ldflags
	Version = "unknown"
	Commit  = "unknown"
	Built   = "unknown"

	// BuildInt is an integer representation of Built
	BuiltInt int
)

func init() {
	// Convert Built into BuiltTime
	var err error
	BuiltInt, err = strconv.Atoi(Built)
	if err != nil {
		panic("unable to convert build-time variable Built into integer")
	}
}
