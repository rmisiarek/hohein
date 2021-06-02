package main

import (
	"io"
)

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

type HohinOptions struct {
	pathHosts   string
	pathHeaders string
	pathValues  string
	output      string
	timeout     int
}

func NewHohin(o *HohinOptions) (*Hohin, error) {
	sourceHosts, err := validateSource(o.pathHosts)
	if err != nil {
		return nil, err
	}

	sourceHeaders, err := readSource(o.pathHeaders)
	if err != nil {
		return nil, err
	}

	sourceValues, err := readSource(o.pathValues)
	if err != nil {
		return nil, err
	}

	return &Hohin{
		sourceHosts:   sourceHosts,
		sourceHeaders: sourceHeaders,
		sourceValues:  sourceValues,
		client:        getClient(o.timeout),
		options:       o,
	}, nil
}

type Hohin struct {
	sourceHeaders []string
	sourceValues  []string
	sourceHosts   io.ReadCloser
	client        *HttpClient
	options       *HohinOptions
}

func (h *Hohin) Start() {

}
