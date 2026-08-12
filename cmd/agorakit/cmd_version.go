package main

import (
	"flag"
	"fmt"

	"github.com/felicabrera/agora/internal/version"
)

func runVersion(args []string) error {
	fs := flag.NewFlagSet("version", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("version: unexpected argument %q", fs.Arg(0))
	}
	fmt.Printf("agorakit %s\n", version.Read())
	return nil
}
