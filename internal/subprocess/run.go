/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package subprocess

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hartwork/go-wait-for-it/internal/logging"
)

func RunCommand(argv []string, log logging.Log) error {
	command := exec.Command(argv[0], argv[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	log.Neutral("Running command: %s", strings.Join(argv, " "))
	err := command.Run()

	if err == nil {
		log.Success("Command succeeded.")
	} else {
		log.Error("Error: %v", err)
	}
	return err
}

func ExitCodeFrom(err error) int {
	if err == nil {
		return 0
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}

	if _, ok := err.(*exec.Error); ok {
		return 127
	}

	return 1
}
