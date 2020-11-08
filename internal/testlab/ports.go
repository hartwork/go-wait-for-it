/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package testlab

import (
	"net"
	"testing"

	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/stretchr/testify/require"
)

func WithAvailablePort(t *testing.T, actUpon func(syntax.Address)) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	defer listener.Close()

	address, err := syntax.ParseAddress(listener.Addr().String())
	require.Nil(t, err)

	actUpon(address)
}

func WithUnavailablePort(t *testing.T, actUpon func(syntax.Address)) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)

	address, err := syntax.ParseAddress(listener.Addr().String())
	require.Nil(t, err)

	listener.Close()

	actUpon(address)
}
