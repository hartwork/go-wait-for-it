/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package syntax

import (
	"fmt"
	"net"
	"strconv"
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

	host, portText, err := net.SplitHostPort(text)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(portText)
	if err != nil {
		return
	}

	if port <= 0 || port >= 1<<16 {
		return
	}

	return Address{host, uint16(port)}, nil
}
