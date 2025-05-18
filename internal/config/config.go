package config

import (
	"flag"
	"fmt"
)

const (
	DEFAULT_FILE = "gotem.config.json"
)

type Flag struct {
	File         string
	RequestName  string
	OverrideUrl  string
	ListRequests bool
}

func InitFlags(args []string) (Flag, error) {
	fs := flag.NewFlagSet("gotem", flag.ContinueOnError)

	fileName := fs.String("f", DEFAULT_FILE, "Config file path")
	reqName := fs.String("req-name", "", "API request name to execute")
	ovrUrl := fs.String("url", "", "Overrides config url")
	listReq := fs.Bool("ls", false, "Display requests available")

	err := fs.Parse(args)
	if err != nil {
		return Flag{}, fmt.Errorf("InitFlags: parse: %w", err)
	}

	newFlags := Flag{
		File:         *fileName,
		RequestName:  *reqName,
		OverrideUrl:  *ovrUrl,
		ListRequests: *listReq,
	}

	return newFlags, nil
}

// NOTE: Idea for flags
// in: request as input, (uses no config file)
// header: flags to override config file requests
