package main

import (
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

func (h *HttpClient) request(payload Payload, debug bool) (*HohinResult, error) {
	request, err := baseRequest(payload)
	if err != nil {
		return nil, err
	}

	request.Header.Set(payload.key, payload.value)

	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	nh, nv := normalizeHeader(response.Header)
	reflectedKey, _ := isHeaderKeyReflected(nh, payload.key)
	reflectedValue, _ := isHeaderValueReflected(nv, payload.value)
	reflectedValueBody := isValueReflectedInBody(response.Body, payload.value)
	location := getLocation(response)

	var confirmed bool
	if inIntSlice(successCodes, response.StatusCode) {
		confirmed, err = h.confirmVulnerability(payload)
		if err != nil {
			return nil, err
		}
	}

	result := &HohinResult{
		payload:            payload,
		responseStatusCode: response.StatusCode,
		responseURL:        response.Request.URL.String(),
		location:           location,
		reflectedKey:       reflectedKey,
		reflectedValue:     reflectedValue,
		reflectedValueBody: reflectedValueBody,
		confirmed:          confirmed,
	}

	return result, nil
}

func baseRequest(payload Payload) (*http.Request, error) {
	req, err := http.NewRequest(payload.method, payload.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Connection", "close")
	req.Close = true

	return req, nil
}

func getLocation(response *http.Response) string {
	location := ""
	_location, err := response.Location()
	if err != nil {
		location = ""
	} else {
		location = _location.String()
	}

	return location
}

func isValueReflectedInBody(response io.Reader, testValue string) string {
	testValue = strings.ToLower(testValue)
	body, err := ioutil.ReadAll(response)
	if err != nil {
		return ""
	}

	if strings.Contains(strings.ToLower(string(body)), testValue) {
		return testValue
	}

	return ""
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

func isHeaderKeyReflected(headers []string, testHeaderKey string) (string, bool) {
	testHeaderKey = strings.ToLower(testHeaderKey)
	for _, header := range headers {
		if header == testHeaderKey {
			return header, true
		}
	}
	return "", false
}

func isHeaderValueReflected(values []string, testHeaderValue string) (string, bool) {
	testHeaderValue = strings.ToLower(testHeaderValue)
	for _, value := range values {
		if value == testHeaderValue {
			return value, true
		}
	}
	return "", false
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
