package main

import (
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
