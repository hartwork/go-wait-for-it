/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/cli"
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
		err = fmt.Errorf("Failed to connected to %s for %s.", address, timeout)
	}

	results <- connectResult{address, duration, err}
}

func waitForMultipleAddressesWithTimeout(addresses []syntax.Address, timeout time.Duration, log logging.Log) (err error) {
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

func runCommand(argv []string, log logging.Log) int {
	command := exec.Command(argv[0], argv[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	log.Neutral("Running command: %s", strings.Join(argv, " "))
	err := command.Run()

	if err == nil {
		log.Success("Command succeeded.")
		return 0
	}

	log.Error("Error: %v", err)

	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}

	return 127
}

func main() {
	config := cli.Parse(os.Args[1:])
	log := logging.Log{Quiet: config.Quiet}

	if err := waitForMultipleAddressesWithTimeout(config.Addresses, config.Timeout, log); err != nil {
		log.Error("Aborting...")
		os.Exit(1)
	}

	if len(config.Argv) > 0 {
		exitCode := runCommand(config.Argv, log)
		os.Exit(exitCode)
	}
}
