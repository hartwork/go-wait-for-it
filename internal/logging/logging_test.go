/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package logging

import (
	"os"
	"testing"

	"github.com/hartwork/go-wait-for-it/internal/mocking"
)

func TestLoggingPlain(t *testing.T) {
	log := Log{Quiet: false}

	mocking.AssertOutputEquals(t, func() {
		log.Neutral("111")
	}, &os.Stdout, "[*] 111\n")

	mocking.AssertOutputEquals(t, func() {
		log.Success("222")
	}, &os.Stdout, "[+] 222\n")

	mocking.AssertOutputEquals(t, func() {
		log.Error("333")
	}, &os.Stdout, "[-] 333\n")
}

func TestLoggingSprintf(t *testing.T) {
	log := Log{Quiet: false}

	mocking.AssertOutputEquals(t, func() {
		log.Neutral("%s %s", "111", "222")
	}, &os.Stdout, "[*] 111 222\n")

	mocking.AssertOutputEquals(t, func() {
		log.Success("%s %s", "333", "444")
	}, &os.Stdout, "[+] 333 444\n")

	mocking.AssertOutputEquals(t, func() {
		log.Error("%s %s", "555", "666")
	}, &os.Stdout, "[-] 555 666\n")
}

func TestLoggingQuiet(t *testing.T) {
	log := Log{Quiet: true}

	mocking.AssertOutputEquals(t, func() {
		log.Neutral("111")
	}, &os.Stdout, "")

	mocking.AssertOutputEquals(t, func() {
		log.Success("222")
	}, &os.Stdout, "")

	mocking.AssertOutputEquals(t, func() {
		log.Error("333")
	}, &os.Stdout, "")
}
