package main

import (
	"io"
	"os"
)

func openStdinOrFile(inputs string) (io.ReadCloser, error) {
	r := os.Stdin

	if inputs != "" {
		var err error

		r, err = openFile(inputs)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func openFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
