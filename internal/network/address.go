/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

type Address interface {
	Host() string
	Port() uint16
	String() string
}

type addressInfo struct {
	host string
	port uint16
}

func NewAddressUnchecked(host string, port uint16) Address {
	return addressInfo{host, port}
}

func (a addressInfo) Host() string {
	return a.host
}

func (a addressInfo) Port() uint16 {
	return a.port
}
