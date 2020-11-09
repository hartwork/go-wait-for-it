/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithListeningPort(t *testing.T) {
	WithListeningPort(t, func(a1 Address) {
		WithListeningPort(t, func(a2 Address) {
			assert.NotEqual(t, 0, a1.Port)
			assert.NotEqual(t, 0, a2.Port)
			assert.NotEqual(t, a1.Port, a2.Port)
		})
	})
}

func TestWithUnusedPort(t *testing.T) {
	WithUnusedPort(t, func(a1 Address) {
		WithUnusedPort(t, func(a2 Address) {
			assert.NotEqual(t, 0, a1.Port)
			assert.NotEqual(t, 0, a2.Port)
			// Note: There is no guarantee that the two ports are unequal
		})
	})
}
