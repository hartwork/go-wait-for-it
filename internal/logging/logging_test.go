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

func assertStdoutEquals(t *testing.T, actUpon func(log Log), expectedOutput string) {
	log := Log{Quiet: false}

	stdoutBackup := os.Stdout
	defer func() {
		os.Stdout = stdoutBackup
	}()

	os.Stdout = newMemoryFile(t)
	defer os.Stdout.Close()

	actUpon(log)

	assert.Equal(t, collect(t, os.Stdout), expectedOutput)
}

func TestLoggingPlain(t *testing.T) {
	assertStdoutEquals(t, func(log Log) {
		log.Neutral("111")
	}, "[*] 111\n")
	assertStdoutEquals(t, func(log Log) {
		log.Success("222")
	}, "[+] 222\n")
	assertStdoutEquals(t, func(log Log) {
		log.Error("333")
	}, "[-] 333\n")
}

func TestLoggingSprintf(t *testing.T) {
	assertStdoutEquals(t, func(log Log) {
		log.Neutral("%s %s", "111", "222")
	}, "[*] 111 222\n")
	assertStdoutEquals(t, func(log Log) {
		log.Success("%s %s", "333", "444")
	}, "[+] 333 444\n")
	assertStdoutEquals(t, func(log Log) {
		log.Error("%s %s", "555", "666")
	}, "[-] 555 666\n")
}
