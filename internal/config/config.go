package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

const (
	DEFAULT_FILE = "gotem.gohttp"
)

type Flag struct {
	File         string
	RequestName  string
	OverrideUrl  string
	ListRequests bool
	ExtraHeaders map[string][]string
}

func InitFlags(args []string) (Flag, error) {
	fs := flag.NewFlagSet("gotem", flag.ContinueOnError)

	fileName := fs.String("f", DEFAULT_FILE, "Config file path")
	reqName := fs.String("req-name", "", "API request name to execute")
	ovrUrl := fs.String("url", "", "Overrides config url")
	listReq := fs.Bool("ls", false, "Display requests available")
	h := headers{}
	fs.Var(&h, "H", "Add header to request, format: <KEY>:<VALUE>")

	err := fs.Parse(args)
	if err != nil {
		return Flag{}, fmt.Errorf("InitFlags: parse: %w", err)
	}

	newFlags := Flag{
		File:         *fileName,
		RequestName:  *reqName,
		OverrideUrl:  *ovrUrl,
		ListRequests: *listReq,
		ExtraHeaders: h,
	}

	return newFlags, nil
}

type headers map[string][]string

func (h *headers) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headers) Set(value string) error {
	if h == nil {
		h = (*headers)(&map[string][]string{})
	}
	if !strings.Contains(value, ":") {
		return errors.New("value is not formatted in key value")
	}
	splitStr := strings.Split(value, ":")
	if len(splitStr) != 2 {
		return errors.New("value containt more than one ':', unable to parse")
	}
	key := splitStr[0]
	val := splitStr[1]

	if len(key) == 0 {
		return errors.New("left side of ':' is empty")
	}

	if len(val) == 0 {
		return errors.New("right side of ':' is empty")
	}

	if h == nil {
		return errors.New("map is nil")
	}

	mapVal, ok := (*h)[key]
	if !ok {
		(*h)[key] = []string{val}
	} else {
		(*h)[key] = append(mapVal, val)
	}

	return nil
}

// NOTE: Idea for flags
// in: request as input, (uses no config file)
