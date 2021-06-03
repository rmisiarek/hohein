package main

import (
	"bufio"
	"fmt"
	"io"
)

type HohinResult struct {
	payloads              Payload
	statusCode            int
	host                  string
	location              string
	reflectedKeys         []string
	reflectedValues       []string
	reflectedValuesInBody []string
	confirmed             bool
}

type Payload struct {
	host      string
	method    string
	payload   map[string]string
	reference int
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
	payloads := h.buildPaylods()
	scanner := bufio.NewScanner(h.sourceHosts)
	for scanner.Scan() {
		host := scanner.Text()
		reference := h.client.referenceStatusCode("GET", host)
		fmt.Println(host)
		fmt.Println(Purple(fmt.Sprintf("\t==> reference status code: %d", reference)))
		for _, payload := range payloads {
			h.client.request(Payload{
				method:    "GET",
				host:      host,
				payload:   payload,
				reference: reference,
			})
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
