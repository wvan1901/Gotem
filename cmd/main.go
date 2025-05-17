package main

import (
	"fmt"
	"github.com/wvan1901/Gotem/internal"
	"github.com/wvan1901/Gotem/internal/config"
	"github.com/wvan1901/Gotem/internal/file"
	"io"
	"os"
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
		return err
	}

	// Get Config File
	reqsFile, err := file.GetConfigFile(flags.File)
	if err != nil {
		return err
	}

	// Get Requests Info
	userReqs, err := internal.ParseJsonInput(reqsFile)
	if err != nil {
		return err
	}

	// List Request
	if flags.ListRequests {
		list := internal.ListRequests(userReqs)
		_, err = w.Write([]byte(list))
		if err != nil {
			return err
		}
		return nil
	}

	// Get Request
	reqs := internal.ConvertInputsToReqs(userReqs)
	req, err := internal.GetRequest(flags.RequestName, reqs)
	if err != nil {
		return err
	}

	// Make A Request
	res, err := req.Execute()
	if err != nil {
		return err
	}

	userRes := internal.ConvertResponse(*res)
	resString, err := userRes.JsonString()
	if err != nil {
		return err
	}

	// Output request
	_, err = w.Write([]byte(resString + "\n"))
	if err != nil {
		return err
	}

	return nil
}
