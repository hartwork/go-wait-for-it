/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package testlab

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTempFile(t *testing.T) *os.File {
	fd, err := ioutil.TempFile("", "wait-for-it-testing-")
	require.Nil(t, err)
	return fd
}

func collect(t *testing.T, file *os.File) string {
	bytes, err := ioutil.ReadFile(file.Name())
	require.Nil(t, err)
	return string(bytes)
}

func WithOutputCapturing(t *testing.T, act func(), fileToMock **os.File) string {
	originalFile := *fileToMock
	defer func() {
		*fileToMock = originalFile
	}()

	*fileToMock = newTempFile(t)
	defer (*fileToMock).Close()
	defer os.Remove((*fileToMock).Name())

	act()

	return collect(t, *fileToMock)
}

func AssertOutputEquals(t *testing.T, act func(), fileToMock **os.File, expectedOutput string) {
	actualOutput := WithOutputCapturing(t, act, fileToMock)
	assert.Equal(t, expectedOutput, actualOutput)
}
