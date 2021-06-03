package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	client *http.Client
}

var successCodes = []int{
	200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
}

func getClient(timeout int) *HttpClient {
	t := time.Duration(timeout) * time.Second

	// transport := &http.Transport{
	// 	MaxIdleConns:      30,
	// 	IdleConnTimeout:   time.Second,
	// 	DisableKeepAlives: true,
	// 	// DisableCompression: true,
	// 	// TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   t,
	// 		KeepAlive: time.Second,
	// 	}).DialContext,
	// }

	client := &http.Client{
		// Transport: transport,
		Timeout: t,
	}

	return &HttpClient{client: client}
}

func (h *HttpClient) referenceStatusCode(method, host string) int {
	request, err := http.NewRequest(method, host, nil)
	if err != nil {
		return 0
	}

	request.Header.Set("Connection", "close")
	request.Close = true

	response, err := h.client.Do(request)
	if err != nil {
		return 0
	}
	defer response.Body.Close()

	return response.StatusCode
}

func (h *HttpClient) confirmVulnerability(payload Payload) (bool, error) {
	request, err := baseRequest(payload)
	if err != nil {
		return false, err
	}

	response, err := h.client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if !inIntSlice(successCodes, response.StatusCode) {
		return true, nil
	}

	return false, nil
}

func (h *HttpClient) request(payloads Payload) (*HohinResult, error) {
	request, err := baseRequest(payloads)
	if err != nil {
		return nil, err
	}

	var header, value string
	for header, value = range payloads.payload {
		request.Header.Set(header, value)
	}

	response, err := h.client.Do(request)
	if err != nil {
		fmt.Println(Red(fmt.Sprintf("\t==> %s", err.Error())))
		return nil, err
	}
	defer response.Body.Close()

	nh, nv := normalizeHeader(response.Header)
	reflectedKeys, found := headerKeysReflected(nh, payloads.payload)
	if found {
		fmt.Println("\t\t>> found reflected headers", reflectedKeys)
		for _, v := range reflectedKeys {
			fmt.Println(Green(fmt.Sprintf("\t\t>> %s", v)))
		}
	}

	reflectedValues, found := headerValuesReflected(nv, payloads.payload)
	if found {
		fmt.Println("\t\t>> found reflected header values")
		for _, v := range reflectedValues {
			fmt.Println(Green(fmt.Sprintf("\t\t>> %s", v)))
		}
	}

	reflectedValuesInBody, found := valuesReflectedInBody(response.Body, payloads.payload)
	if found {
		fmt.Println("\t\t>> found reflected value in body")
		for _, v := range reflectedValuesInBody {
			fmt.Println(Green(fmt.Sprintf("\t\t>> %s", v)))
		}
	}

	location := getLocation(response)

	var confirmed bool
	if inIntSlice(successCodes, response.StatusCode) {
		fmt.Println(Yellow(fmt.Sprintf("\t==> status code: %d | payload: %s | confirming...", response.StatusCode, value)))

		confirmed, err = h.confirmVulnerability(payloads)
		if err != nil {
			fmt.Println(Red(fmt.Sprintf("\t\t>> %s", err.Error())))
			return nil, err
		}

		if confirmed {
			fmt.Println(Green("\t\t>> confirmation: OK!"))
		} else {
			fmt.Println(Red("\t\t>> confirmation: NOK"))
		}
	} else {
		fmt.Println(Blue(fmt.Sprintf("\t==> status code: %d | payload: %s", response.StatusCode, value)))
	}

	result := &HohinResult{
		payloads:              payloads,
		statusCode:            response.StatusCode,
		host:                  response.Request.URL.String(),
		location:              location,
		reflectedKeys:         reflectedKeys,
		reflectedValues:       reflectedValues,
		reflectedValuesInBody: reflectedValuesInBody,
		confirmed:             confirmed,
	}

	return result, nil
}

func baseRequest(payload Payload) (*http.Request, error) {
	req, err := http.NewRequest(payload.method, payload.host, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Connection", "close")
	req.Close = true

	return req, nil
}

func getLocation(response *http.Response) string {
	var location string

	_location, err := response.Location()
	if err != nil {
		location = ""
	} else {
		location = _location.String()
	}

	return location
}

func normalizeHeader(response http.Header) ([]string, []string) {
	headers := []string{}
	values := []string{}

	for k, vals := range response {
		headers = append(headers, strings.ToLower(k))
		for _, v := range vals {
			values = append(values, strings.ToLower(v))
		}
	}

	return headers, values
}

func valuesReflectedInBody(response io.Reader, payloads map[string]string) ([]string, bool) {
	results := []string{}
	var found bool

	body, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, false
	}

	for _, payloadValue := range payloads {
		if strings.Contains(strings.ToLower(string(body)), strings.ToLower(payloadValue)) {
			found = true
			results = append(results, payloadValue)
		}
	}

	return results, found
}

func headerKeysReflected(headers []string, payloads map[string]string) ([]string, bool) {
	results := []string{}
	var found bool

	for payloadHeader := range payloads {
		for _, header := range headers {
			if header == strings.ToLower(payloadHeader) {
				found = true
				results = append(results, header)
			}
		}
	}

	return results, found
}

func headerValuesReflected(values []string, payloads map[string]string) ([]string, bool) {
	results := []string{}
	var found bool

	for _, payloadValue := range payloads {
		for _, value := range values {
			if value == strings.ToLower(payloadValue) {
				found = true
				results = append(results, value)
			}
		}
	}

	return results, found
}

func inIntSlice(s []int, v int) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}

	return false
}

func inStrSlice(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}

	return false
}
