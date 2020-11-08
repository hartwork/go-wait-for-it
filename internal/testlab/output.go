/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package testlab

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func newMemoryFile(t *testing.T) *os.File {
	fd, err := unix.MemfdCreate("fake stdout", 0)
	require.Nil(t, err)
	filename := fmt.Sprintf("/proc/self/fd/%d", fd)
	return os.NewFile(uintptr(fd), filename)
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

	*fileToMock = newMemoryFile(t)
	defer (*fileToMock).Close()

	act()

	return collect(t, *fileToMock)
}

func AssertOutputEquals(t *testing.T, act func(), fileToMock **os.File, expectedOutput string) {
	actualOutput := WithOutputCapturing(t, act, fileToMock)
	assert.Equal(t, expectedOutput, actualOutput)
}
