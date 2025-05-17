package file

import (
	"fmt"
	"io"
	"os"
)

func GetConfigFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("GetJsonFile: os: %w", err)
	}
	var fileErr error
	defer func() {
		fileErr = file.Close()
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("GetJsonFile: io: %w", err)
	}

	return fileBytes, fileErr
}
