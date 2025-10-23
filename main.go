package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wvan1901/Gotem/internal/cli"
	"github.com/wvan1901/Gotem/internal/config"
	"github.com/wvan1901/Gotem/internal/file"
)

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	flags, err := config.InitFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	// Get Config File
	reqsFile, err := file.GetConfigFile(flags.File)
	if err != nil {
		return err
	}

	// Get Program Info aka Requests
	prog, err := cli.ParseInputIntoProgram(string(reqsFile))
	if err != nil {
		return err
	}

	// List Request
	if flags.ListRequests {
		list := cli.ListAstRequests(prog.Requests)
		_, err = w.Write([]byte(list))
		if err != nil {
			return err
		}
		return nil
	}

	// Get single request
	req, err := cli.GetAstRequest(flags.RequestName, prog.Requests)
	if err != nil {
		return err
	}

	// TODO: Look into Override/Add Custom Labels via flags

	// Make a request
	resp, err := cli.MakeRequest(req, flags.ExtraHeaders)
	if err != nil {
		return err
	}

	// Turn resp into a json string
	respJsonBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// Output request
	_, err = w.Write(append(respJsonBytes, '\n'))
	if err != nil {
		return err
	}

	return nil
}
