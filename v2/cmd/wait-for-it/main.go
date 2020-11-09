/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package main

import (
	"os"

	"github.com/hartwork/go-wait-for-it/v2/internal/cli"
	"github.com/hartwork/go-wait-for-it/v2/internal/logging"
	"github.com/hartwork/go-wait-for-it/v2/internal/network"
	"github.com/hartwork/go-wait-for-it/v2/internal/subprocess"
)

var exitFunc func(int) = os.Exit

func innerMain(argv []string) error {
	config, err := cli.Parse(argv[1:])
	if err != nil {
		return err
	}
	if config == nil { // i.e. help output has just been presented
		return nil
	}

	log := map[bool]func() logging.Log{
		true:  logging.NewNullLog,
		false: logging.NewStdoutLog,
	}[config.Quiet]()

	if err := network.WaitForMultipleAddressesWithTimeout(config.Addresses, config.Timeout, log); err != nil {
		log.Error("Aborting...")
		return err
	}

	if len(config.Argv) > 0 {
		return subprocess.RunCommand(config.Argv, log)
	}

	return nil
}

func main() {
	exitFunc(subprocess.ExitCodeFrom(innerMain(os.Args)))
}
