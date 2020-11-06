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

func (l Log) log(prefix string, format string, args []interface{}) {
	if l.Quiet {
		return
	}

	var message string

	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	fmt.Printf("%s %s\n", prefix, message)
}

func (l Log) Neutral(format string, args ...interface{}) {
	l.log("[*]", format, args)
}

func (l Log) Error(format string, args ...interface{}) {
	l.log("[-]", format, args)
}

func (l Log) Success(format string, args ...interface{}) {
	l.log("[+]", format, args)
}
