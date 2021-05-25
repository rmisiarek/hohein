package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHeaderKeyReflected(t *testing.T) {
	headers := []string{"header-1", "header-2"}

	header, ok := isHeaderKeyReflected(headers, "header-1")
	assert.Equal(t, "header-1", header)
	assert.Equal(t, true, ok)

	header, ok = isHeaderKeyReflected(headers, "header-3")
	assert.Equal(t, "", header)
	assert.Equal(t, false, ok)
}

func TestIsHeaderValueReflected(t *testing.T) {
	values := []string{"value-1", "value-2"}

	value, ok := isHeaderValueReflected(values, "value-1")
	assert.Equal(t, "value-1", value)
	assert.Equal(t, true, ok)

	value, ok = isHeaderValueReflected(values, "value-3")
	assert.Equal(t, "", value)
	assert.Equal(t, false, ok)
}

func TestNormalizeHeader(t *testing.T) {
	H := map[string][]string{
		"HEADER-1": {"value-1"},
		"header-2": {"VALUE-2"},
	}

	headers, values := normalizeHeader(H)
	assert.Equal(t, true, inStrSlice(headers, "header-1"))
	assert.Equal(t, true, inStrSlice(headers, "header-2"))
	assert.Equal(t, true, inStrSlice(values, "value-1"))
	assert.Equal(t, true, inStrSlice(values, "value-2"))
}

func TestIsValueReflectedInBody(t *testing.T) {
	r := ioutil.NopCloser(strings.NewReader("short text just for testing"))
	value := isValueReflectedInBody(r, "JUST for")
	assert.Equal(t, "just for", value)

	value = isValueReflectedInBody(r, "will be not found")
	assert.Equal(t, "", value)

	var err errorReader
	value = isValueReflectedInBody(err, "throw error")
	assert.Equal(t, "", value)
}

type errorReader int

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func TestGetLocation(t *testing.T) {
	responseOK := http.Response{
		Header: http.Header{"Location": []string{"example.com"}},
	}
	responseNOK := http.Response{
		Header: http.Header{"Location": []string{""}},
	}

	location := getLocation(&responseOK)
	assert.Equal(t, "example.com", location)

	location = getLocation(&responseNOK)
	assert.Equal(t, "", location)
}

func TestBaseRequest(t *testing.T) {
	payloadOK := Payload{
		url:    "example.com",
		method: "GET",
	}
	payloadNOK := Payload{
		url:    "example.com",
		method: "()", // to throw error, as validMethod() == false
	}

	request, err := baseRequest(payloadOK)
	assert.NotNil(t, request)
	assert.Nil(t, err)
	assert.Equal(t, "close", request.Header.Get("Connection"))
	assert.Equal(t, true, request.Close)

	request, err = baseRequest(payloadNOK)
	assert.Nil(t, request)
	assert.NotNil(t, err)
}
