// Package integration tests the omar CLI end-to-end.
//
// These tests build the omar binary and invoke it as a subprocess, testing
// argument parsing, output formatting, and error handling with real system
// dependencies (or their absence).
//
// Run via:
//
//	go test -tags=integration ./tests/integration/ -v
//
// Or from the project root:
//
//	make test-integration
package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// omarBin is the path to the built omar binary, set in TestMain.
var omarBin string

// projectRoot returns the project root by walking up from the test directory
// until it finds go.mod.
func projectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Fprintf(os.Stderr, "could not find go.mod from %s\n", dir)
			os.Exit(1)
		}
		dir = parent
	}
}

func TestMain(m *testing.M) {
	// Build the omar binary before running any tests.
	bin := filepath.Join(os.TempDir(), "omar-integration-test")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	root := projectRoot()
	cmd := exec.Command("go", "build",
		"-ldflags=-X 'github.com/nevotheless/omar/internal/version.Version=0.0.0-integration-test'",
		"-o", bin,
		"./cmd/omar",
	)
	cmd.Dir = root
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build omar: %v\n", err)
		os.Exit(1)
	}
	omarBin = bin

	// Ensure cleanup after all tests.
	defer os.Remove(bin)

	os.Exit(m.Run())
}

// execOmar runs the omar binary with the given arguments and returns stdout,
// stderr, and the exit error.
func execOmar(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(omarBin, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// ─── version ──────────────────────────────────────────────────────────────

func TestVersion_Text(t *testing.T) {
	stdout, stderr, err := execOmar("version")
	require.NoError(t, err, "omar version should succeed")
	assert.Empty(t, stderr, "stderr should be empty")
	assert.Contains(t, stdout, "omar CLI version")
	assert.Contains(t, stdout, "0.0.0-integration-test")
}

func TestVersion_JSON(t *testing.T) {
	stdout, stderr, err := execOmar("version", "--json")
	require.NoError(t, err, "omar version --json should succeed")
	assert.Empty(t, stderr, "stderr should be empty")

	var v struct {
		CLIVersion      string `json:"cli_version"`
		ImageVersion    string `json:"image_version"`
		SystemImage     string `json:"system_image"`
	}
	require.True(t, json.Valid([]byte(stdout)), "output must be valid JSON")
	require.NoError(t, json.Unmarshal([]byte(stdout), &v))
	assert.Equal(t, "0.0.0-integration-test", v.CLIVersion)
	// image_version and system_image depend on the host (may be empty on non-ostree)
	t.Logf("version --json: cli=%s image=%s system=%s", v.CLIVersion, v.ImageVersion, v.SystemImage)
}

// ─── status ───────────────────────────────────────────────────────────────

func TestStatus(t *testing.T) {
	stdout, stderr, err := execOmar("status")
	require.NoError(t, err, "omar status should succeed")
	assert.Empty(t, stderr, "stderr should be empty")
	assert.Contains(t, stdout, "Image:")
	assert.Contains(t, stdout, "Version:")
	assert.Contains(t, stdout, "Booted:")
	assert.Contains(t, stdout, "Staged:")
	assert.Contains(t, stdout, "Rollback:")
}

func TestStatus_OutputFormat(t *testing.T) {
	stdout, _, err := execOmar("status")
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	require.GreaterOrEqual(t, len(lines), 5, "status output should have at least 5 lines")

	// First line should be the header
	assert.Equal(t, "=== omar System Status ===", lines[0],
		"first line should be the status header, got: %q", lines[0])
	// Second line should be "  Image: ..."
	assert.True(t, strings.HasPrefix(lines[1], "  Image:"),
		"second line should start with '  Image:', got: %q", lines[1])
	// Third line should be "  Version: ..."
	assert.True(t, strings.HasPrefix(lines[2], "  Version:"),
		"third line should start with '  Version:', got: %q", lines[2])
}

// ─── help ─────────────────────────────────────────────────────────────────

func TestHelp(t *testing.T) {
	stdout, stderr, err := execOmar("--help")
	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "omar")
	assert.Contains(t, stdout, "Usage:")
	assert.Contains(t, stdout, "Available Commands:")
	assert.Contains(t, stdout, "install")
	assert.Contains(t, stdout, "status")
	assert.Contains(t, stdout, "update")
	assert.Contains(t, stdout, "rollback")
	assert.Contains(t, stdout, "pkg")
	assert.Contains(t, stdout, "version")
}

func TestHelpShort(t *testing.T) {
	stdout, stderr, err := execOmar("-h")
	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "omar")
}

// ─── error handling ───────────────────────────────────────────────────────

func TestUnknownCommand(t *testing.T) {
	_, stderr, err := execOmar("nonexistent")
	assert.Error(t, err, "unknown command should fail")
	assert.Contains(t, stderr, "unknown command")
	assert.Contains(t, stderr, "nonexistent")
}

func TestVersion_InvalidFlag(t *testing.T) {
	_, stderr, err := execOmar("version", "--invalid")
	assert.Error(t, err)
	assert.Contains(t, stderr, "unknown flag")
}

// ─── pkg ──────────────────────────────────────────────────────────────────

func TestPkgAdd_NoArgs(t *testing.T) {
	_, stderr, err := execOmar("pkg", "add")
	assert.Error(t, err, "pkg add without args should fail")
	assert.Contains(t, stderr, "accepts 1 arg")
}

func TestPkgAdd_Help(t *testing.T) {
	stdout, stderr, err := execOmar("pkg", "add", "--help")
	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Contains(t, stdout, "omar pkg add")
}

func TestPkgList_NoBackends(t *testing.T) {
	// pkg list should work even without flatpak/distrobox (returns empty)
	stdout, stderr, err := execOmar("pkg", "list")
	require.NoError(t, err, "pkg list should not fail")
	assert.Empty(t, stderr)
	// Should return empty or some output listing nothing
	// It's okay for this to be empty or contain an error message
	t.Logf("pkg list output: %q", stdout)
}

// ─── install ──────────────────────────────────────────────────────────────

func TestInstallNoArgs(t *testing.T) {
	// omar install without flags should produce some output (help or error)
	stdout, stderr, err := execOmar("install")
	t.Logf("install stdout: %q", stdout)
	t.Logf("install stderr: %q", stderr)
	// Either shows help or tries to convert (will fail gracefully on non-ostree)
	if err != nil {
		assert.Contains(t, stderr, "bootc") // likely bootc-related error
	}
}

func TestInstallFreshNoDisk(t *testing.T) {
	_, stderr, err := execOmar("install", "--fresh")
	assert.Error(t, err, "--fresh without --disk should fail")
	assert.Contains(t, stderr, "disk")
}

// ─── update ───────────────────────────────────────────────────────────────

func TestUpdateNoOstree(t *testing.T) {
	// On a non-ostree system, omar update should fail or report no immutable OS
	_, stderr, err := execOmar("update")
	if err != nil {
		// Should mention ostree or bootc or registry
		t.Logf("update stderr: %q", stderr)
	}
}

// ─── rollback ─────────────────────────────────────────────────────────────

func TestRollbackNoOstree(t *testing.T) {
	_, stderr, err := execOmar("rollback")
	t.Logf("rollback err=%v stderr=%q", err, stderr)
	if err != nil {
		assert.Contains(t, stderr, "not a bootc system")
	}
}

// ─── version compatibility ────────────────────────────────────────────────

func TestVersionJSON_StableFields(t *testing.T) {
	// Ensure the JSON output has stable field names (no breaking changes).
	// cli_version and booted are always present.
	// image, version, error are conditionally present (omitempty).
	stdout, _, err := execOmar("version", "--json")
	require.NoError(t, err)

	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(stdout), &raw))

	// Always-present fields
	assert.Contains(t, raw, "cli_version", "JSON must contain cli_version")
	assert.Contains(t, raw, "booted", "JSON must contain booted")
	assert.Equal(t, "0.0.0-integration-test", raw["cli_version"])

	// On a mutable system, image should be present (with "none" value)
	if img, ok := raw["image"]; ok {
		assert.Contains(t, img.(string), "none")
	}
}

// keysOf returns the keys of a map for diagnostic purposes.
func keysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
