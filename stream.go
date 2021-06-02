package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func readSource(src string) ([]string, error) {
	if src != "" && !fileExists(src) {
		return nil, fmt.Errorf("%s does not exist", src)
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}

func validateSource(src string) (io.ReadCloser, error) {
	if src != "" && !fileExists(src) {
		return nil, fmt.Errorf("%s does not exist", src)
	}

	source, err := openStdinOrFile(src)
	if err != nil {
		return nil, err
	}

	return source, nil
}

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
