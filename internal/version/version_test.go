package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion_Default(t *testing.T) {
	// Default value when not set via ldflags
	assert.Equal(t, "0.0.0-dev", Version)
}

func TestVersion_SetViaLdflags(t *testing.T) {
	// Simulate what ldflags does at build time
	old := Version
	Version = "2026.5.0"
	defer func() { Version = old }()

	assert.Equal(t, "2026.5.0", Version)
}

func TestVersion_Empty(t *testing.T) {
	// Edge case: should never be empty in production,
	// but test behavior if it were
	old := Version
	Version = ""
	defer func() { Version = old }()

	assert.Empty(t, Version)
}
