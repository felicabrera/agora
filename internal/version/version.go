// Package version reports what this binary is.
//
// The values come from the Go toolchain's own VCS stamping rather than from
// -ldflags, so they cannot be set to something the build did not actually come
// from. For a system whose central claim is verifiability, a binary that can
// state which commit produced it is the least it can do.
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Info describes the running binary.
type Info struct {
	Revision  string
	Modified  bool
	GoVersion string
}

// Read returns the build information stamped into this binary.
func Read() Info {
	info := Info{Revision: "unknown", GoVersion: runtime.Version()}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			info.Revision = s.Value
		case "vcs.modified":
			info.Modified = s.Value == "true"
		}
	}
	return info
}

// String renders the build information for --version output and logs.
func (i Info) String() string {
	rev := i.Revision
	if i.Modified {
		rev += "-dirty"
	}
	return fmt.Sprintf("%s (%s)", rev, i.GoVersion)
}
