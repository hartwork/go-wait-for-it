/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package logging

import (
	"fmt"
)

type Log interface {
	Neutral(format string, args ...interface{})
	Error(format string, args ...interface{})
	Success(format string, args ...interface{})
}

type stdoutLog struct{}

type nullLog struct{}

func NewStdoutLog() Log {
	return stdoutLog{}
}

func NewNullLog() Log {
	return nullLog{}
}

func (l stdoutLog) log(prefix string, format string, args []interface{}) {
	var message string

	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	fmt.Printf("%s %s\n", prefix, message)
}

func (l stdoutLog) Neutral(format string, args ...interface{}) {
	l.log("[*]", format, args)
}

func (l nullLog) Neutral(format string, args ...interface{}) {
	// no-op
}

func (l stdoutLog) Error(format string, args ...interface{}) {
	l.log("[-]", format, args)
}

func (l nullLog) Error(format string, args ...interface{}) {
	// no-op
}

func (l stdoutLog) Success(format string, args ...interface{}) {
	l.log("[+]", format, args)
}

func (l nullLog) Success(format string, args ...interface{}) {
	// no-op
}
