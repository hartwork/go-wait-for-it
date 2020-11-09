/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressToString(t *testing.T) {
	require.Equal(t,
		"hostname:123",
		NewAddressUnchecked("hostname", 123).String(),
	)
}

func TestErrorToString(t *testing.T) {
	require.Equal(t,
		"Malformed address: no port here",
		MalformedAddressError{"no port here"}.Error(),
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
		assert.Equal(t, test.expectedHost, address.Host())
		assert.Equal(t, test.expectedPort, address.Port())
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
		assert.Equal(t, "Malformed address: "+candidate, err.Error())
	}
}
