package bootc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusJSON_Parse(t *testing.T) {
	// Valid bootc status JSON (as produced by bootc 1.1+)
	input := `{
		"booted": {
			"id": "deployment-abc123",
			"image": {
				"image": "ghcr.io/nevotheless/omar:rolling",
				"version": "2026.5.0"
			},
			"timestamp": 1715000000
		},
		"staged": null,
		"rollback": null
	}`

	var s Status
	err := json.Unmarshal([]byte(input), &s)
	require.NoError(t, err)

	require.NotNil(t, s.Booted)
	assert.Equal(t, "deployment-abc123", s.Booted.ID)
	assert.Equal(t, "ghcr.io/nevotheless/omar:rolling", s.Booted.Image.Image)
	assert.Equal(t, "2026.5.0", s.Booted.Image.Version)
	assert.Equal(t, int64(1715000000), s.Booted.Timestamp)
	assert.Nil(t, s.Staged)
	assert.Nil(t, s.Rollback)
}

func TestStatusJSON_ParseWithStaged(t *testing.T) {
	// When an update is staged
	input := `{
		"booted": {
			"id": "booted-001",
			"image": {
				"image": "ghcr.io/nevotheless/omar:rolling",
				"version": "2026.5.0"
			},
			"timestamp": 1715000000
		},
		"staged": {
			"id": "staged-002",
			"image": {
				"image": "ghcr.io/nevotheless/omar:rolling",
				"version": "2026.6.0"
			},
			"timestamp": 1715100000
		}
	}`

	var s Status
	err := json.Unmarshal([]byte(input), &s)
	require.NoError(t, err)

	require.NotNil(t, s.Booted)
	require.NotNil(t, s.Staged)
	assert.Equal(t, "2026.6.0", s.Staged.Image.Version)
}

func TestStatusJSON_ParseMinimal(t *testing.T) {
	// Minimal valid JSON — booted may be null on fresh installs
	input := `{"booted": null}`

	var s Status
	err := json.Unmarshal([]byte(input), &s)
	require.NoError(t, err)
	assert.Nil(t, s.Booted)
}

func TestStatusJSON_ParseInvalidJSON(t *testing.T) {
	// Invalid JSON should error
	input := `{invalid}`

	var s Status
	err := json.Unmarshal([]byte(input), &s)
	assert.Error(t, err)
}

func TestImageRef_MissingFields(t *testing.T) {
	// Version should be empty when not provided
	input := `{
		"booted": {
			"id": "test-id",
			"image": {
				"image": "test-image"
			},
			"timestamp": 0
		}
	}`

	var s Status
	err := json.Unmarshal([]byte(input), &s)
	require.NoError(t, err)
	require.NotNil(t, s.Booted)
	assert.Equal(t, "test-image", s.Booted.Image.Image)
	assert.Empty(t, s.Booted.Image.Version)
}

func TestDeployment_JSONTags(t *testing.T) {
	// Verify the JSON struct tags match bootc output format
	d := Deployment{
		ID: "test-deploy",
		Image: ImageRef{
			Image:   "my-image:tag",
			Version: "1.0.0",
		},
		Timestamp: 1234567890,
	}

	b, err := json.Marshal(d)
	require.NoError(t, err)

	var decoded Deployment
	err = json.Unmarshal(b, &decoded)
	require.NoError(t, err)
	assert.Equal(t, d.ID, decoded.ID)
	assert.Equal(t, d.Image.Image, decoded.Image.Image)
	assert.Equal(t, d.Image.Version, decoded.Image.Version)
}
