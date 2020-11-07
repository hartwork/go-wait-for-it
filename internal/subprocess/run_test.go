/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package subprocess

import (
	"fmt"
	"os"
	"testing"

	"github.com/hartwork/go-wait-for-it/internal/logging"
	"github.com/hartwork/go-wait-for-it/internal/mocking"
	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	log := logging.Log{Quiet: true}
	mocking.AssertOutputEquals(t, func() {
		mocking.AssertOutputEquals(t, func() {
			RunCommand([]string{"sh", "-c", "echo 111; echo 222 >&2"}, log)
		}, &os.Stderr, "222\n")
	}, &os.Stdout, "111\n")
}

func TestExitCodeFromRegular(t *testing.T) {
	tests := []struct {
		argv         []string
		expectedCode int
	}{
		{[]string{"true"}, 0},
		{[]string{"sh", "-c", "exit 123"}, 123},
		{[]string{"b1788316d2acd022536d5b750e6eb6af01e2ca0e"}, 127}, // no such command
	}

	log := logging.Log{Quiet: true}
	for _, test := range tests {
		actualCode := ExitCodeFrom(RunCommand(test.argv, log))
		assert.Equal(t, test.expectedCode, actualCode)
	}

	assert.Equal(t, 1, ExitCodeFrom(fmt.Errorf("Some arbitrary error")))
}
