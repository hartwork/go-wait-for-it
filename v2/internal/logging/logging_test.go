/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package logging

import (
	"os"
	"testing"

	"github.com/hartwork/go-wait-for-it/v2/internal/testlab"
)

func TestLoggingPlain(t *testing.T) {
	log := NewStdoutLog()

	testlab.AssertOutputEquals(t, func() {
		log.Neutral("111")
	}, &os.Stdout, "[*] 111\n")

	testlab.AssertOutputEquals(t, func() {
		log.Success("222")
	}, &os.Stdout, "[+] 222\n")

	testlab.AssertOutputEquals(t, func() {
		log.Error("333")
	}, &os.Stdout, "[-] 333\n")
}

func TestLoggingSprintf(t *testing.T) {
	log := NewStdoutLog()

	testlab.AssertOutputEquals(t, func() {
		log.Neutral("%s %s", "111", "222")
	}, &os.Stdout, "[*] 111 222\n")

	testlab.AssertOutputEquals(t, func() {
		log.Success("%s %s", "333", "444")
	}, &os.Stdout, "[+] 333 444\n")

	testlab.AssertOutputEquals(t, func() {
		log.Error("%s %s", "555", "666")
	}, &os.Stdout, "[-] 555 666\n")
}

func TestLoggingQuiet(t *testing.T) {
	log := NewNullLog()

	testlab.AssertOutputEquals(t, func() {
		log.Neutral("111")
	}, &os.Stdout, "")

	testlab.AssertOutputEquals(t, func() {
		log.Success("222")
	}, &os.Stdout, "")

	testlab.AssertOutputEquals(t, func() {
		log.Error("333")
	}, &os.Stdout, "")
}
