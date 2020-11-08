/*
 * Copyright (C) 2020 Sebastian Pipping <sebastian@pipping.org>
 * Licensed under AGPL v3 or later
 */
package testlab

import (
	"testing"

	"github.com/hartwork/go-wait-for-it/internal/syntax"
	"github.com/stretchr/testify/assert"
)

func TestWithAvailablePort(t *testing.T) {
	WithAvailablePort(t, func(a1 syntax.Address) {
		WithAvailablePort(t, func(a2 syntax.Address) {
			assert.NotEqual(t, a1.Port, 0)
			assert.NotEqual(t, a2.Port, 0)
			assert.NotEqual(t, a1.Port, a2.Port)
		})
	})
}

func TestWithUnavailablePort(t *testing.T) {
	WithUnavailablePort(t, func(a1 syntax.Address) {
		WithUnavailablePort(t, func(a2 syntax.Address) {
			assert.NotEqual(t, a1.Port, 0)
			assert.NotEqual(t, a2.Port, 0)
			// Note: There is no guarantee that the two ports are unequal
		})
	})
}
