/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"testing"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestWaitForAddressSuccess(t *testing.T) {
	WithListeningPort(t, func(address Address) {
		port := address.Port()

		addresses := []Address{
			NewAddressUnchecked("localhost", port),
			NewAddressUnchecked("127.0.0.1", port),
			NewAddressUnchecked("[::]", port),
			NewAddressUnchecked("", port),
		}

		deadlineCombined := time.After(2 * time.Second)
		for _, address := range addresses {
			available := waitForAddressForever(address)
			select {
			case <-available:
			case <-deadlineCombined:
				t.Errorf("waitForAddressForever should be long done by now.")
			}
		}
	})
}

func TestWaitForAddressFailure(t *testing.T) {
	WithUnusedPort(t, func(address Address) {
		timeout := 1250 * time.Millisecond // small to not blow up test runtime
		select {
		case <-waitForAddressForever(address):
			t.Errorf("waitForAddressForever was expected to never finish.")
		case <-time.After(timeout):
		}
	})
}

func TestWaitForAddressWithTimeoutSuccess(t *testing.T) {
	WithListeningPort(t, func(address Address) {
		timeout := 2 * time.Second
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, results)

		result := <-results
		assert.Nil(t, result.err)
	})
}

func TestWaitForAddressWithTimeoutFailure(t *testing.T) {
	WithUnusedPort(t, func(address Address) {
		timeout := 100 * time.Millisecond // small to not blow up test runtime
		results := make(chan connectResult)

		go waitForAddressWithTimeout(address, timeout, results)

		result := <-results
		assert.NotNil(t, result.err)
	})
}

func TestWaitForMultipleAddressesWithTimeoutSuccess(t *testing.T) {
	WithListeningPort(t, func(a1 Address) {
		WithListeningPort(t, func(a2 Address) {
			addresses := []Address{a1, a2}
			timeout := 2 * time.Second
			log := logging.NewNullLog()

			err := WaitForMultipleAddressesWithTimeout(addresses, timeout, log)

			assert.Nil(t, err)
		})
	})
}

func TestWaitForMultipleAddressesWithTimeoutFailure(t *testing.T) {
	WithListeningPort(t, func(a1 Address) {
		WithUnusedPort(t, func(a2 Address) {
			addresses := []Address{a1, a2}
			timeout := 100 * time.Millisecond // small to not blow up test runtime
			log := logging.NewNullLog()

			err := WaitForMultipleAddressesWithTimeout(addresses, timeout, log)

			assert.NotNil(t, err)
		})
	})
}
