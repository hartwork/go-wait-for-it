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

	"github.com/hartwork/go-wait-for-it/cli"
	"github.com/hartwork/go-wait-for-it/logging"
	"github.com/hartwork/go-wait-for-it/syntax"
)

type ConnectResult struct {
	address  syntax.Address
	duration time.Duration
	success  bool
}

func waitForAddress(address syntax.Address) chan bool {
	available := make(chan bool)
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

func waitForAddressWithTimeout(address syntax.Address, timeout time.Duration, results chan ConnectResult) {
	duration := timeout
	success := false

	deadline := make(<-chan time.Time)
	if timeout > 0 {
		deadline = time.After(timeout)
	}

	before := time.Now()

	select {
	case <-waitForAddress(address):
		duration = time.Now().Sub(before)
		success = true
	case <-deadline:
	}

	results <- ConnectResult{address, duration, success}
}

func waitForMultipleAddressesWithTimeout(addresses []syntax.Address, timeout time.Duration, log logging.Log) bool {
	success := true
	results := make(chan ConnectResult)

	for _, address := range addresses {
		log.Neutral(fmt.Sprintf("Trying to connect to %s...", address))
		go waitForAddressWithTimeout(address, timeout, results)
	}

	for replies := 0; success && replies < len(addresses); replies++ {
		result := <-results
		if result.success {
			log.Success(fmt.Sprintf("Connected to %s after %s.", result.address, result.duration))
		} else {
			success = false
			log.Error(fmt.Sprintf("Failed to connected to %s for %s.", result.address, result.duration))
			log.Error("Aborting...")
		}
	}

	return success
}

func runCommand(argv []string, log logging.Log) int {
	command := exec.Command(argv[0], argv[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	log.Neutral(fmt.Sprintf("Running command: %s", strings.Join(argv, " ")))
	err := command.Run()

	if err == nil {
		log.Success("Command succeeded.")
		return 0
	}

	log.Error(fmt.Sprintf("Error: %v", err))

	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}

	return 127
}

func main() {
	config := cli.Parse(os.Args)
	log := logging.Log{Quiet: config.Quiet}

	if !waitForMultipleAddressesWithTimeout(config.Addresses, config.Timeout, log) {
		os.Exit(1)
	}

	if len(config.Argv) > 0 {
		exitCode := runCommand(config.Argv, log)
		os.Exit(exitCode)
	}
}
