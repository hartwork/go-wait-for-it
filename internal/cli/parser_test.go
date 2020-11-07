/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package cli

import (
	"testing"
	"time"

	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/stretchr/testify/assert"
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
				[]syntax.Address{{"h1", 1}, {"h2", 2}},
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
		config := Parse(test.args)
		assert.Equal(t, config, test.expectedConfig)
	}
}
