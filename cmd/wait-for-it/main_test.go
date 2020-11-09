/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package main

import (
	"os"
	"testing"

	"github.com/hartwork/go-wait-for-it/internal/network"
	"github.com/hartwork/go-wait-for-it/internal/testlab"
	"github.com/stretchr/testify/assert"
)

func TestInnerMainUsageError(t *testing.T) {
	testlab.WithOutputCapturing(t, func() {
		err := innerMain([]string{"wait-for-it", "--no-such-thing"})
		assert.NotNil(t, err)
	}, &os.Stderr)
}

func TestInnerMainHelpDisplayed(t *testing.T) {
	testlab.WithOutputCapturing(t, func() {
		err := innerMain([]string{"wait-for-it", "--help"})
		assert.Nil(t, err)
	}, &os.Stdout)
}

func TestInnerMainTimout(t *testing.T) {
	network.WithUnusedPort(t, func(address network.Address) {
		err := innerMain([]string{"wait-for-it", "-q", "-t", "1", "-s", address.String()})
		assert.NotNil(t, err)
	})
}

func TestInnerMainRunError(t *testing.T) {
	err := innerMain([]string{"wait-for-it", "-q", "false"})
	assert.NotNil(t, err)
}

func TestInnerMainSuccessWithCommand(t *testing.T) {
	err := innerMain([]string{"wait-for-it", "-q", "true"})
	assert.Nil(t, err)
}

func TestInnerMainSuccessWithoutCommand(t *testing.T) {
	err := innerMain([]string{"wait-for-it"})
	assert.Nil(t, err)
}

func TestMain(t *testing.T) {
	originalOsArgs := os.Args
	originalExitFunc := exitFunc

	os.Args = []string{"wait-for-it", "-q", "--", "sh", "-c", "exit 123"}
	exitFunc = func(code int) {
		assert.Equal(t, 123, code)
	}

	defer func() {
		os.Args = originalOsArgs
		exitFunc = originalExitFunc
	}()

	main()
}
