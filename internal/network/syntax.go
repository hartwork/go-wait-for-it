/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func (a addressInfo) String() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}

type MalformedAddressError struct {
	value string
}

func (e MalformedAddressError) Error() string {
	return fmt.Sprintf("Malformed address: %s", e.value)
}

func ParseAddress(text string) (address Address, networkError error) {
	networkError = MalformedAddressError{text}

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

	return addressInfo{host, uint16(port)}, nil
}
