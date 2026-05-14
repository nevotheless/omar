package image

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIVersion(t *testing.T) {
	v := CLIVersion()
	assert.NotEmpty(t, v, "CLIVersion() should return a non-empty string")
	assert.Equal(t, "0.0.0-dev", v)
}

func TestString_NoOstree(t *testing.T) {
	// When there is no ostree, Current() returns a "none" info
	s := String()
	assert.Contains(t, s, "none (mutable system)")
	assert.Contains(t, s, "Booted:      false")
	assert.Contains(t, s, "Staged:      false")
	assert.Contains(t, s, "Rollback:    false")
}

func TestCurrent_NoOstree(t *testing.T) {
	info, err := Current()
	require.NoError(t, err)
	require.NotNil(t, info)
	assert.Equal(t, "none (mutable system)", info.Image)
	assert.Empty(t, info.Version)
	assert.Empty(t, info.DeploymentID)
	assert.False(t, info.Booted)
	assert.False(t, info.Staged)
	assert.False(t, info.RollbackExists)
}

func TestString_Format(t *testing.T) {
	// Verify the output format contains expected labels
	s := String()
	assert.True(t, strings.HasPrefix(s, "Image:"), "output should start with Image:")
	assert.Contains(t, s, "Version:")
}
