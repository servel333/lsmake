package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || errors.Is(err, os.ErrExist)
}

func resolveIncludedFilePath(baseFile, includeFile string) string {
	baseDir := filepath.Dir(baseFile)
	includedPath, err := filepath.Abs(filepath.Join(baseDir, includeFile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		return filepath.Join(baseDir, includeFile)
	}
	return includedPath
}
