/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"fmt"
	"net"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/logging"
	"github.com/hartwork/go-wait-for-it/internal/syntax"
)

type connectResult struct {
	address  syntax.Address
	duration time.Duration
	err      error
}

func waitForAddress(address syntax.Address) <-chan bool {
	available := make(chan bool, 1)
	go func() {
		for {
			c, err := net.Dial("tcp", address.String())
			if err == nil {
				available <- true
				c.Close()
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return available
}

func waitForAddressWithTimeout(address syntax.Address, timeout time.Duration, startedAt time.Time, results chan<- connectResult) {
	duration := timeout

	deadline := make(<-chan time.Time)
	if timeout > 0 {
		deadline = time.After(timeout)
	}

	err := error(nil)

	select {
	case <-waitForAddress(address):
		duration = time.Now().Sub(startedAt)
	case <-deadline:
		err = fmt.Errorf("Failed to connect to %s for %s.", address, timeout)
	}

	results <- connectResult{address, duration, err}
}

func WaitForMultipleAddressesWithTimeout(addresses []syntax.Address, timeout time.Duration, log logging.Log) (err error) {
	results := make(chan connectResult, len(addresses))
	startedAt := time.Now()

	for _, address := range addresses {
		log.Neutral("Trying to connect to %s...", address)
		go waitForAddressWithTimeout(address, timeout, startedAt, results)
	}

	for range addresses {
		if result := <-results; result.err == nil {
			log.Success("Connected to %s after %s.", result.address, result.duration)
		} else {
			log.Error(result.err.Error())
			err = result.err // the first error is as good as the last, here
		}
	}
	close(results)

	return err
}
