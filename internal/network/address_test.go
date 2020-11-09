/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddressUnchecked(t *testing.T) {
	address := NewAddressUnchecked("h1", 123)
	assert.Equal(t, "h1", address.Host())
	assert.Equal(t, uint16(123), address.Port())
}
