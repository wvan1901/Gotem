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
	flags := config.InitFlags(args)
	err := flags.IsValid()
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
	fmt.Println("Wicho: res", res.StatusCode, string(res.Body), req.Name)

	userRes := internal.ConvertResponse(*res)
	output := fmt.Sprint("", userRes)
	// Output request to stdout
	w.Write([]byte(string(output) + "\n"))

	return nil
}
