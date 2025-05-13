package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Stdin, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(r io.Reader, w io.Writer, args []string) error {
	// TODO: flags

	return nil
}
