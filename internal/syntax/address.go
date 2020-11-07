/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package syntax

import (
	"fmt"
	"net"
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

	host, portText, err := net.SplitHostPort(text)
	if err != nil {
		return
	}

	if len(portText) > 0 && portText[0] == '0' { // deny leading zeros
		return
	}

	port, err := strconv.Atoi(portText)
	if err != nil {
		return
	}

	if port <= 0 || port >= 1<<16 {
		return
	}

	if strings.Contains(host, ":") {
		host = "[" + host + "]" // wrapping of IPv6 addresses
	}

	return Address{host, uint16(port)}, nil
}
