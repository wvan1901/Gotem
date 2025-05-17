package main

import (
	"fmt"
	"gotem/internal"
	"gotem/internal/config"
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
	reqsFile, err := internal.GetJsonFile(flags.File)
	if err != nil {
		return err
	}

	// Get requests info
	userReq, err := internal.ParseJsonInput(reqsFile)
	if err != nil {
		return err
	}
	reqs := internal.ConvertInputsToReqs(userReq)

	// Get Request
	req, err := internal.GetRequest(flags.RequestName, reqs)
	if err != nil {
		return err
	}
	// Make a request
	res, err := req.Execute()
	if err != nil {
		return err
	}

	userRes := internal.ConvertResponse(*res)
	resString, err := userRes.JsonString()
	if err != nil {
		return err
	}

	// Output request to stdout
	_, err = w.Write([]byte(resString + "\n"))
	if err != nil {
		return err
	}

	return nil
}
