/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func WithListeningPort(t *testing.T, actUpon func(Address)) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	defer listener.Close()

	address, err := ParseAddress(listener.Addr().String())
	require.Nil(t, err)

	actUpon(address)
}

func WithUnusedPort(t *testing.T, actUpon func(Address)) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)

	address, err := ParseAddress(listener.Addr().String())
	require.Nil(t, err)

	listener.Close()

	actUpon(address)
}
