package main

import (
	"fmt"
	"io"
)

type Hohin struct {
	source  io.ReadCloser
	client  *HttpClient
	options *HohinOptions
}

type HohinOptions struct {
	source  string
	output  string
	timeout int
}

type HohinResult struct {
	payload            Payload
	responseStatusCode int
	responseURL        string
	reflectedKey       string
	reflectedValue     string
	reflectedValueBody string
	location           string
	confirmed          bool
}

type Payload struct {
	url    string
	method string
	key    string
	value  string
}

func NewHohin(o *HohinOptions) (*Hohin, error) {
	if o.source != "" && !fileExists(o.source) {
		return nil, fmt.Errorf("%s does not exist", o.source)
	}

	source, err := openStdinOrFile(o.source)
	if err != nil {
		return nil, fmt.Errorf("error while opening %s", o.source)
	}

	return &Hohin{
		source:  source,
		client:  getClient(o.timeout),
		options: o,
	}, nil
}
