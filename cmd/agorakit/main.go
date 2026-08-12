// Command agorakit is the ÁGORA operator CLI.
//
// It follows ÁBACO's shape: stdlib flag only, one file per subcommand, manual
// dispatch in main. No CLI framework, because the dependency surface of a voting
// system is part of its threat model and a flag parser is not worth widening it.
//
// Today it carries the two subcommands that are useful before any protocol code
// exists: `version`, and `keygen`, which exercises the ported crypto core end to
// end and gives the team something real to check the build against. Election
// setup, ballot casting and tally arrive with the MVP.
package main

import (
	"flag"
	"fmt"
	"os"
)

// commands is the dispatch table. Adding a subcommand means adding a file with a
// run function and one line here.
var commands = map[string]struct {
	summary string
	run     func(args []string) error
}{
	"version": {"print build information", runVersion},
	"keygen":  {"generate a threshold election key", runKeygen},
}

func main() {
	flag.Usage = usage
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	name := os.Args[1]
	if name == "-h" || name == "--help" || name == "help" {
		usage()
		return
	}
	cmd, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "agorakit: unknown command %q\n\n", name)
		usage()
		os.Exit(2)
	}
	if err := cmd.run(os.Args[2:]); err != nil {
		fatalf("%v", err)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, "agorakit — ÁGORA operator CLI\n\nusage: agorakit <command> [flags]\n\ncommands:\n")
	for _, name := range []string{"version", "keygen"} {
		fmt.Fprintf(os.Stderr, "  %-10s %s\n", name, commands[name].summary)
	}
	fmt.Fprint(os.Stderr, "\nRun 'agorakit <command> -h' for the flags of a command.\n")
}

// fatalf reports a fatal error and exits, matching ÁBACO's convention of a
// single exit helper rather than log.Fatal scattered around.
func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "agorakit: "+format+"\n", args...)
	os.Exit(1)
}
