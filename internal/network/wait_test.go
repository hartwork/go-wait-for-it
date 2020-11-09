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

type addressMock struct {
	succeed bool
}

func (a addressMock) Host() string {
	return "none"
}

func (a addressMock) Port() uint16 {
	return 0
}

func (a addressMock) String() string {
	return "none:0"
}

func (a addressMock) waitForForever() <-chan bool {
	available := make(chan bool, 1)
	if a.succeed {
		available <- true
	}
	return available
}

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
			available := address.waitForForever()
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
		case <-address.waitForForever():
			t.Errorf("waitForAddressForever was expected to never finish.")
		case <-time.After(timeout):
		}
	})
}

func TestWaitForAddressWithTimeoutSuccess(t *testing.T) {
	address := addressMock{succeed: true}
	timeout := 2 * time.Second
	results := make(chan connectResult)

	go waitForAddressWithTimeout(address, timeout, results)

	result := <-results
	assert.Nil(t, result.err)
}

func TestWaitForAddressWithTimeoutFailure(t *testing.T) {
	address := addressMock{succeed: false}
	timeout := 100 * time.Millisecond // small to not blow up test runtime
	results := make(chan connectResult)

	go waitForAddressWithTimeout(address, timeout, results)

	result := <-results
	assert.NotNil(t, result.err)
}

func TestWaitForMultipleAddressesWithTimeoutSuccess(t *testing.T) {
	addresses := []Address{
		addressMock{succeed: true},
		addressMock{succeed: true},
	}
	timeout := 2 * time.Second
	log := logging.NewNullLog()

	err := WaitForMultipleAddressesWithTimeout(addresses, timeout, log)

	assert.Nil(t, err)
}

func TestWaitForMultipleAddressesWithTimeoutFailure(t *testing.T) {
	addresses := []Address{
		addressMock{succeed: true},
		addressMock{succeed: false},
	}
	timeout := 100 * time.Millisecond // small to not blow up test runtime
	log := logging.NewNullLog()

	err := WaitForMultipleAddressesWithTimeout(addresses, timeout, log)

	assert.NotNil(t, err)
}
