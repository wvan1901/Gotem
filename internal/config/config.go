package config

import (
	"errors"
	"flag"
	"fmt"
)

const (
	DEFAULT_FILE = "gotem.config.json"
)

type Flag struct {
	File        string
	RequestName string
}

func InitFlags(args []string) (Flag, error) {
	fs := flag.NewFlagSet("gotem", flag.ContinueOnError)

	fileName := fs.String("f", DEFAULT_FILE, "config file name")
	reqName := fs.String("req-name", "", "API request name to execute")

	err := fs.Parse(args)
	if err != nil {
		return Flag{}, fmt.Errorf("InitFlags: parse: %w", err)
	}

	newFlags := Flag{
		File:        *fileName,
		RequestName: *reqName,
	}

	err = newFlags.IsValid()
	if err != nil {
		return Flag{}, fmt.Errorf("InitFlags: validate: %w", err)
	}
	return newFlags, nil
}

func (f *Flag) IsValid() error {
	if f.File == "" {
		return errors.New("file name can't be empty")
	}
	if f.RequestName == "" {
		return errors.New("request name can't be empty")
	}
	return nil
}

// NOTE: Idea for flags
// ls: list all the requests with a description
// in: request as input, (uses no config file)
