/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"testing"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/hartwork/go-wait-for-it/internal/testlab"
	"github.com/stretchr/testify/assert"
)

func TestWaitForAddress(t *testing.T) {
	testlab.WithAvailablePort(t, func(address syntax.Address) {
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
	testlab.WithAvailablePort(t, func(address syntax.Address) {
		timeout := 2 * time.Second
		startedAt := time.Now()
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, startedAt, results)

		result := <-results
		assert.Nil(t, result.err)
	})
}

func TestWaitForAddressWithTimeoutFailure(t *testing.T) {
	testlab.WithUnavailablePort(t, func(address syntax.Address) {
		timeout := 100 * time.Millisecond // small to not blow up test runtime
		startedAt := time.Now()
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, startedAt, results)

		result := <-results
		assert.NotNil(t, result.err)
	})
}
