package main

import (
	"bufio"
	"io"
)

type Payload struct {
	host      string
	method    string
	headers   map[string]string
	reference int
}

type HohinOptions struct {
	pathHosts   string
	pathHeaders string
	pathValues  string
	output      string
	timeout     int
}

func NewHohin(options *HohinOptions) (*Hohin, error) {
	sourceHosts, err := validateSource(options.pathHosts)
	if err != nil {
		return nil, err
	}

	sourceHeaders, err := readSource(options.pathHeaders)
	if err != nil {
		return nil, err
	}

	sourceValues, err := readSource(options.pathValues)
	if err != nil {
		return nil, err
	}

	return &Hohin{
		sourceHosts:   sourceHosts,
		sourceHeaders: sourceHeaders,
		sourceValues:  sourceValues,
		client:        getClient(options.timeout),
		options:       options,
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
	payloads := h.buildPaylods()
	scanner := bufio.NewScanner(h.sourceHosts)
	for scanner.Scan() {
		host := scanner.Text()
		reference := h.client.referenceStatusCode("GET", host)

		whiteBold.Printf("%s\n", host)
		purple.Printf("\t==> reference status code: %d\n", reference)

		out := make(chan ResultWrapper)
		for _, payload := range payloads {
			go func(payload map[string]string) {
				out <- h.client.request(Payload{
					method:    "GET",
					host:      host,
					headers:   payload,
					reference: reference,
				})
			}(payload)
		}

		for i := 0; i < len(payloads); i++ {
			handleResult(<-out)
		}
	}
}

type payloads []map[string]string

func (h *Hohin) buildPaylods() payloads {
	results := make(payloads, len(h.sourceValues))

	for i, value := range h.sourceValues {
		m := make(map[string]string)
		for _, header := range h.sourceHeaders {
			m[header] = value
		}
		results[i] = m
	}

	return results
}
