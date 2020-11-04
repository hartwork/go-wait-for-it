/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package logging

import (
	"fmt"
)

type Log struct {
	Quiet bool
}

func (l Log) log(prefix, message string) {
	if l.Quiet {
		return
	}
	fmt.Printf("%s %s\n", prefix, message)
}

func (l Log) Neutral(message string) {
	l.log("[*]", message)
}

func (l Log) Error(message string) {
	l.log("[-]", message)
}

func (l Log) Success(message string) {
	l.log("[+]", message)
}
