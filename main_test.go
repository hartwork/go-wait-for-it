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
	"github.com/stretchr/testify/require"
)

func TestWaitForAddress(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	require.Nil(t, err)
	defer listener.Close()

	address, err := syntax.ParseAddress(listener.Addr().String())
	require.Nil(t, err)
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
}
