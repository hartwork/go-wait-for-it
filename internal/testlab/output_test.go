/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package testlab

import (
	"fmt"
	"os"
	"testing"
)

func TestAssertOutputEquals(t *testing.T) {
	AssertOutputEquals(t, func() {
		fmt.Fprintln(os.Stderr, "hello test helper")
	}, &os.Stderr, "hello test helper\n")
}
