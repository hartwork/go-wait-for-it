/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"fmt"
	"net"
	"time"

	"github.com/hartwork/go-wait-for-it/v2/internal/logging"
)

type connectResult struct {
	address   Address
	stoppedAt time.Time
	err       error
}

func (a addressInfo) waitForForever() <-chan bool {
	available := make(chan bool, 1)
	go func() {
		for {
			c, err := net.DialTimeout("tcp", a.String(), 500*time.Millisecond)
			if err == nil {
				available <- true
				c.Close()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return available
}

func waitForAddressWithTimeout(address Address, timeout time.Duration, results chan<- connectResult) {
	deadline := make(<-chan time.Time)
	if timeout > 0 {
		deadline = time.After(timeout)
	}

	err := error(nil)

	select {
	case <-address.waitForForever():
	case <-deadline:
		err = fmt.Errorf("Failed to connect to %s for %s.", address, timeout)
	}

	results <- connectResult{address, time.Now(), err}
}

func WaitForMultipleAddressesWithTimeout(addresses []Address, timeout time.Duration, log logging.Log) (err error) {
	results := make(chan connectResult, len(addresses))
	startedAt := time.Now()

	for _, address := range addresses {
		log.Neutral("Trying to connect to %s...", address)
		go waitForAddressWithTimeout(address, timeout, results)
	}

	for range addresses {
		if result := <-results; result.err == nil {
			duration := result.stoppedAt.Sub(startedAt)
			log.Success("Connected to %s after %s.", result.address, duration)
		} else {
			log.Error(result.err.Error())
			err = result.err // the first error is as good as the last, here
		}
	}
	close(results)

	return err
}
