/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package main

import (
	"net"
	"testing"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withAvailablePort(t *testing.T, actUpon func(syntax.Address)) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	defer listener.Close()

	address, err := syntax.ParseAddress(listener.Addr().String())
	require.Nil(t, err)

	actUpon(address)
}

func withUnavailablePort(t *testing.T, actUpon func(syntax.Address)) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)

	address, err := syntax.ParseAddress(listener.Addr().String())
	require.Nil(t, err)

	listener.Close()

	actUpon(address)
}

func TestWaitForAddress(t *testing.T) {
	withAvailablePort(t, func(address syntax.Address) {
		port := address.Port

		addresses := []syntax.Address{
			{"localhost", port},
			{"127.0.0.1", port},
			{"[::]", port},
			{"", port},
		}

		deadlineCombined := time.After(2 * time.Second)
		for _, address := range addresses {
			available := waitForAddress(address)
			select {
			case <-available:
			case <-deadlineCombined:
				t.Errorf("waitForAddress should be long done by now.")
			}
		}
	})
}

func TestWaitForAddressWithTimeoutSuccess(t *testing.T) {
	withAvailablePort(t, func(address syntax.Address) {
		timeout := 2 * time.Second
		startedAt := time.Now()
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, startedAt, results)

		result := <-results
		assert.Nil(t, result.err)
	})
}

func TestWaitForAddressWithTimeoutFailure(t *testing.T) {
	withUnavailablePort(t, func(address syntax.Address) {
		timeout := 100 * time.Millisecond // small to not blow up test runtime
		startedAt := time.Now()
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, startedAt, results)

		result := <-results
		assert.NotNil(t, result.err)
	})
}
