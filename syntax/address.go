/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package syntax

import (
	"fmt"
	"strconv"
	"strings"
)

type Address struct {
	host string
	port uint16
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}

type MalformedAddressError struct {
	value string
}

func (e MalformedAddressError) Error() string {
	return fmt.Sprintf("Malformed address: %s", e.value)
}

func ParseAddress(text string) (address Address, syntaxError error) {
	syntaxError = MalformedAddressError{text}

	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	if port <= 0 || port >= 1<<16 {
		return
	}

	return Address{parts[0], uint16(port)}, nil
}
