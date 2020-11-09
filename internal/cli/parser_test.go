/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package cli

import (
	"os"
	"testing"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/network"
	"github.com/hartwork/go-wait-for-it/internal/testlab"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserWellFormed(t *testing.T) {
	tests := []struct {
		args           []string
		expectedConfig Config
	}{
		{ // all default settings
			[]string{},
			Config{
				nil,
				[]string{},
				false,
				15 * time.Second,
			},
		}, { // cover as many non-defaults combined as possible
			[]string{
				"-q",
				"-t", "2",
				"-s", "h1:1",
				"--service", "h2:2",
				"--",
				"echo", "hello",
			}, Config{
				[]network.Address{
					network.NewAddressUnchecked("h1", 1),
					network.NewAddressUnchecked("h2", 2),
				},
				[]string{"echo", "hello"},
				true,
				2 * time.Second,
			},
		}, { // cover remaining non-defaults
			[]string{
				"--quiet",
				"--timeout", "2",
				"echo", "hello",
			}, Config{
				nil,
				[]string{"echo", "hello"},
				true,
				2 * time.Second,
			},
		},
	}

	for _, test := range tests {
		config, err := Parse(test.args)
		assert.Nil(t, err)
		require.NotNil(t, config)
		assert.Equal(t, test.expectedConfig, *config)
	}
}

func TestParserHelpOutput(t *testing.T) {
	expectedOutput := dedent.Dedent(`
		Wait for service(s) to be available before executing a command.

		Usage:
		  wait-for-it [flags] [-s|--service [HOST]:PORT]... [--] [COMMAND [ARG ..]]

		Flags:
		  -h, --help              help for wait-for-it
		  -q, --quiet             do not output any status messages
		  -s, --service strings   services to test (format '[HOST]:PORT')
		  -t, --timeout uint      timeout in seconds, 0 for no timeout (default 15)
		  -v, --version           version for wait-for-it

		go-wait-for-it is software libre, licensed under the AGPL v3 or later license.
		Please report bugs at https://github.com/hartwork/go-wait-for-it/issues.  Thank you!
	`)[1:] // drop leading newline

	testlab.AssertOutputEquals(t, func() {
		config, err := Parse([]string{"--help"})
		assert.Nil(t, err)
		assert.Nil(t, config)
	}, &os.Stdout, expectedOutput)
}

func TestParserVersionOutput(t *testing.T) {
	testlab.AssertOutputEquals(t, func() {
		config, err := Parse([]string{"--version"})
		assert.Nil(t, err)
		assert.Nil(t, config)
	}, &os.Stdout, "wait-for-it 1.0.0\n")
}

func TestParserUnknownFlag(t *testing.T) {
	stderrOutput := testlab.WithOutputCapturing(t, func() {
		config, err := Parse([]string{"--no-such-thing"})
		assert.NotNil(t, err)
		assert.Nil(t, config)
	}, &os.Stderr)
	assert.Contains(t, stderrOutput, "unknown flag: --no-such-thing")
}

func TestParserMalformedAddress(t *testing.T) {
	stderrOutput := testlab.WithOutputCapturing(t, func() {
		config, err := Parse([]string{"-s", "no colon here"})
		assert.NotNil(t, err)
		assert.Nil(t, config)
	}, &os.Stderr)
	assert.Contains(t, stderrOutput, "Malformed address: no colon here")
}
