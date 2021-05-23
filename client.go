package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	client *http.Client
}

const successCodes = []int{
	200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
}

func inIntSlice(s []int, v int) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
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

func (h *HttpClient) confirmVulnerabilityRequest(potentialVuln *HohinResult) (interface{}, error) {
	request, err := http.NewRequest("GET", potentialVuln.url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Connection", "close")
	request.Close = true

	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if !inIntSlice(successCodes, response.StatusCode) {
		potentialVuln.confirmed = true
	}

	return nil, nil
}

func (h *HttpClient) request(url, method, headerKey, headerValue string) (*HohinResult, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// referenceResponse, err := h.client.Get(url)
	// if err != nil {
	// 	return nil, err
	// }
	// defer referenceResponse.Body.Close()
	// fmt.Printf("referenceResponse %d\n", referenceResponse.ContentLength)

	// req.Header.Add(headerKey, headerValue)
	// req.Header.Add("Connection", "close")
	// req.Close = true

	response, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// fmt.Printf("%#v\n", response.Header)

	// fmt.Printf("resp %d\n", response.ContentLength)

	// ref, val, ok := reflected(response.Header, "via", "google", false)
	// fmt.Printf("%s - %s - %v \n", ref, val, ok)

	nh, nv := normalizeHeader(response.Header, true)
	isHeaderKeyReflected(nh, "via", true)
	isHeaderValueReflected(nv, "1.1 google", true)

	result := &HohinResult{
		statusCode:  response.StatusCode,
		url:         response.Request.URL.String(),
		headerKey:   headerKey,
		headerValue: headerValue,
	}

	return result, nil
}

func normalizeHeader(response http.Header, debug bool) ([]string, []string) {
	headers := []string{}
	values := []string{}

	for k, vals := range response {
		headers = append(headers, strings.ToLower(k))
		for _, v := range vals {
			if debug {
				fmt.Println(strings.ToLower(k), ":", strings.ToLower(v))
			}
			values = append(values, strings.ToLower(v))
		}
	}

	return headers, values
}

func isHeaderKeyReflected(headers []string, testHeaderKey string, debug bool) (string, bool) {
	testHeaderKey = strings.ToLower(testHeaderKey)
	for _, header := range headers {
		if header == testHeaderKey {
			if debug {
				fmt.Printf("found key reflected: %s\n", header)
			}
			return header, true
		}
	}
	return "", false
}

func isHeaderValueReflected(values []string, testHeaderValue string, debug bool) (string, bool) {
	testHeaderValue = strings.ToLower(testHeaderValue)
	for _, value := range values {
		if value == testHeaderValue {
			if debug {
				fmt.Printf("found value reflected: %s\n", value)
			}
			return value, true
		}
	}
	return "", false
}
