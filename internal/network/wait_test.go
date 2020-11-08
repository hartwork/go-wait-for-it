/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"testing"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/logging"
	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/hartwork/go-wait-for-it/internal/testlab"
	"github.com/stretchr/testify/assert"
)

func TestWaitForAddressSuccess(t *testing.T) {
	testlab.WithListeningPort(t, func(address syntax.Address) {
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

func TestWaitForAddressFailure(t *testing.T) {
	testlab.WithUnusedPort(t, func(address syntax.Address) {
		timeout := 1250 * time.Millisecond // small to not blow up test runtime
		select {
		case <-waitForAddress(address):
			t.Errorf("waitForAddress was expected to never finish.")
		case <-time.After(timeout):
		}
	})
}

func TestWaitForAddressWithTimeoutSuccess(t *testing.T) {
	testlab.WithListeningPort(t, func(address syntax.Address) {
		timeout := 2 * time.Second
		startedAt := time.Now()
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, startedAt, results)

		result := <-results
		assert.Nil(t, result.err)
	})
}

func TestWaitForAddressWithTimeoutFailure(t *testing.T) {
	testlab.WithUnusedPort(t, func(address syntax.Address) {
		timeout := 100 * time.Millisecond // small to not blow up test runtime
		startedAt := time.Now()
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, startedAt, results)

		result := <-results
		assert.NotNil(t, result.err)
	})
}

func TestWaitForMultipleAddressesWithTimeoutSuccess(t *testing.T) {
	testlab.WithListeningPort(t, func(a1 syntax.Address) {
		testlab.WithListeningPort(t, func(a2 syntax.Address) {
			addresses := []syntax.Address{a1, a2}
			timeout := 2 * time.Second
			log := logging.Log{Quiet: true}

			err := WaitForMultipleAddressesWithTimeout(addresses, timeout, log)

			assert.Nil(t, err)
		})
	})
}

func TestWaitForMultipleAddressesWithTimeoutFailure(t *testing.T) {
	testlab.WithListeningPort(t, func(a1 syntax.Address) {
		testlab.WithUnusedPort(t, func(a2 syntax.Address) {
			addresses := []syntax.Address{a1, a2}
			timeout := 100 * time.Millisecond // small to not blow up test runtime
			log := logging.Log{Quiet: true}

			err := WaitForMultipleAddressesWithTimeout(addresses, timeout, log)

			assert.NotNil(t, err)
		})
	})
}
