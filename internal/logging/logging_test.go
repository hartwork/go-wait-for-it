/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package logging

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

func assertOutputEquals(t *testing.T, act func(), fileToMock **os.File, expectedOutput string) {
	originalFile := *fileToMock
	defer func() {
		*fileToMock = originalFile
	}()

	*fileToMock = newMemoryFile(t)
	defer (*fileToMock).Close()

	act()

	assert.Equal(t, collect(t, *fileToMock), expectedOutput)
}

func TestLoggingPlain(t *testing.T) {
	log := Log{Quiet: false}

	assertOutputEquals(t, func() {
		log.Neutral("111")
	}, &os.Stdout, "[*] 111\n")

	assertOutputEquals(t, func() {
		log.Success("222")
	}, &os.Stdout, "[+] 222\n")

	assertOutputEquals(t, func() {
		log.Error("333")
	}, &os.Stdout, "[-] 333\n")
}

func TestLoggingSprintf(t *testing.T) {
	log := Log{Quiet: false}

	assertOutputEquals(t, func() {
		log.Neutral("%s %s", "111", "222")
	}, &os.Stdout, "[*] 111 222\n")

	assertOutputEquals(t, func() {
		log.Success("%s %s", "333", "444")
	}, &os.Stdout, "[+] 333 444\n")

	assertOutputEquals(t, func() {
		log.Error("%s %s", "555", "666")
	}, &os.Stdout, "[-] 555 666\n")
}
