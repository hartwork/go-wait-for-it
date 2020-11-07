/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressToString(t *testing.T) {
	require.Equal(t,
		Address{"hostname", 123}.String(),
		"hostname:123",
	)
}

func TestErrorToString(t *testing.T) {
	require.Equal(t,
		MalformedAddressError{"no port here"}.Error(),
		"Malformed address: no port here",
	)
}

func TestParserWellFormed(t *testing.T) {
	tests := []struct {
		candidate    string
		expectedHost string
		expectedPort uint16
	}{
		{"host:1", "host", 1},         // minimum port number
		{"host:65535", "host", 65535}, // maximum port number
		{":123", "", 123},             // hostname absent
		{"[::]:123", "[::]", 123},     // wrapped IPv6 address
	}
	for _, test := range tests {
		address, err := ParseAddress(test.candidate)
		assert.Nil(t, err)
		assert.Equal(t, address.Host, test.expectedHost)
		assert.Equal(t, address.Port, test.expectedPort)
	}
}

func TestParserMalformed(t *testing.T) {
	tests := []string{
		"host",       // no port A
		"host:",      // no port B
		"host:foo",   // port not an integer
		"host:0",     // below minimum port number
		"host:65536", // above maximum port number
		"host:01",    // forbidden leading zeros
		":::1",       // IPv6 address without [..] wrapper
	}

	for _, candidate := range tests {
		_, err := ParseAddress(candidate)
		assert.Equal(t, err.Error(), "Malformed address: "+candidate)
	}
}