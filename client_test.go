package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_headerKeysReflected(t *testing.T) {
	responseHeaders := []string{"header-1", "header-2"}
	payloadsOK := map[string]string{
		"header-1": "value-1",
		"header-2": "value-2",
		"header-3": "value-3",
	}
	payloadsNOK := map[string]string{
		"header-4": "value-4",
		"header-5": "value-5",
		"header-6": "value-6",
	}

	reflectedHeaders, ok := headerKeysReflected(responseHeaders, payloadsOK)
	assert.Equal(t, true, inStrSlice(reflectedHeaders, "header-1"))
	assert.Equal(t, true, inStrSlice(reflectedHeaders, "header-2"))
	assert.Equal(t, false, inStrSlice(reflectedHeaders, "header-3"))
	assert.Equal(t, true, ok)

	reflectedHeaders, ok = headerKeysReflected(responseHeaders, payloadsNOK)
	assert.Equal(t, false, inStrSlice(reflectedHeaders, "header-4"))
	assert.Equal(t, false, inStrSlice(reflectedHeaders, "header-5"))
	assert.Equal(t, false, inStrSlice(reflectedHeaders, "header-6"))
	assert.Equal(t, false, ok)
}

func Test_headerValuesReflected(t *testing.T) {
	responseValues := []string{"value-1", "value-2"}
	payloadsOK := map[string]string{
		"header-1": "value-1",
		"header-2": "value-2",
		"header-3": "value-3",
	}
	payloadsNOK := map[string]string{
		"header-4": "value-4",
		"header-5": "value-5",
		"header-6": "value-6",
	}

	reflectedValues, ok := headerValuesReflected(responseValues, payloadsOK)
	assert.Equal(t, true, inStrSlice(reflectedValues, "value-1"))
	assert.Equal(t, true, inStrSlice(reflectedValues, "value-2"))
	assert.Equal(t, false, inStrSlice(reflectedValues, "value-3"))
	assert.Equal(t, true, ok)

	reflectedValues, ok = headerValuesReflected(responseValues, payloadsNOK)
	assert.Equal(t, false, inStrSlice(reflectedValues, "value-4"))
	assert.Equal(t, false, inStrSlice(reflectedValues, "value-5"))
	assert.Equal(t, false, inStrSlice(reflectedValues, "value-6"))
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

func Test_valuesReflectedInBody(t *testing.T) {
	r := ioutil.NopCloser(strings.NewReader("short body text just for testing"))

	reflectedValues, found := valuesReflectedInBody(r, map[string]string{"header-1": "JUST for", "header-2": "not"})
	assert.Equal(t, true, inStrSlice(reflectedValues, "JUST for"))
	assert.Equal(t, false, inStrSlice(reflectedValues, "head"))
	assert.Equal(t, true, found)

	reflectedValues, found = valuesReflectedInBody(r, map[string]string{"header-1": "will not be found"})
	assert.Equal(t, false, inStrSlice(reflectedValues, "value to found"))
	assert.Equal(t, false, found)

	var err errorReader
	_, found = valuesReflectedInBody(err, map[string]string{"header-1": "error"})
	assert.Equal(t, false, found)
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
		host:   "example.com",
		method: "GET",
	}
	payloadNOK := Payload{
		host:   "example.com",
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
