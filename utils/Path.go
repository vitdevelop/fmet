package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func IsPath(path string) bool {
	return strings.Contains(path, "/")
}

func CurrentPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s\n", err))
	}
	return path, nil
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
