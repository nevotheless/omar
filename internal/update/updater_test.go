package update

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

// Test helpers for mocking exec.Command
type fakeCmd struct {
	output []byte
	err    error
}

func fakeExecCommand(output string, err error) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cmd := exec.Command("echo", output)
		if err != nil {
			cmd = exec.Command("false")
		}
		return cmd
	}
}

func TestCheckRegistry_BootcUpgradeCheckAvailable(t *testing.T) {
	// Mock bootc upgrade --check reporting an update
	result := bootcUpgradeCheckResult{
		UpdateAvailable: true,
		Image:           "ghcr.io/nevotheless/omar:rolling",
		Version:         "2026.6.0",
	}
	b, _ := json.Marshal(result)

	oldExec := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "bootc" && len(args) == 2 && args[0] == "upgrade" && args[1] == "--check" {
			return exec.Command("echo", string(b))
		}
		return exec.Command("false")
	}
	defer func() { execCommand = oldExec }()

	ver, ok, err := checkRegistry("ghcr.io/nevotheless/omar:rolling")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected update available")
	}
	if ver != "2026.6.0" {
		t.Fatalf("expected version 2026.6.0, got %s", ver)
	}
}

func TestCheckRegistry_NoUpdateAvailable(t *testing.T) {
	// Mock bootc upgrade --check saying no update
	oldExec := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		if name == "bootc" && len(args) == 2 && args[0] == "upgrade" && args[1] == "--check" {
			return exec.Command("echo", "No update available")
		}
		return exec.Command("false")
	}
	defer func() { execCommand = oldExec }()

	_, ok, err := checkRegistry("ghcr.io/nevotheless/omar:rolling")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected no update available")
	}
}

func TestCheckRegistry_bootcFailsFallbackToSkopeo(t *testing.T) {
	// bootc upgrade --check fails (exit 1), skopeo succeeds
	callCount := 0
	oldExec := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		callCount++
		if name == "bootc" && len(args) == 2 && args[0] == "upgrade" && args[1] == "--check" {
			// bootc --check returns error (unsupported)
			return exec.Command("false")
		}
		if name == "skopeo" && args[0] == "inspect" {
			// Return a manifest with version label
			manifest := `{"config":{"labels":{"version":"2026.7.0"}}}`
			return exec.Command("echo", manifest)
		}
		return exec.Command("false")
	}
	defer func() { execCommand = oldExec }()

	// Mock execLookPath so hasSkopeo returns true
	oldLookPath := execLookPath
	execLookPath = func(name string) (string, error) {
		if name == "skopeo" {
			return "/usr/bin/skopeo", nil
		}
		return "", exec.ErrNotFound
	}
	defer func() { execLookPath = oldLookPath }()

	ver, ok, err := checkRegistry("ghcr.io/nevotheless/omar:rolling")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected update available via skopeo")
	}
	if ver != "2026.7.0" {
		t.Fatalf("expected version 2026.7.0, got %s", ver)
	}
}

func TestCheck_NoOstree(t *testing.T) {
	_, err := Check()
	if err == nil {
		t.Fatal("expected error on non-ostree system")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("not a bootc system")) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormalizeImageRef(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"ghcr.io/nevotheless/omar:rolling", "ghcr.io/nevotheless/omar:rolling"},
		{"docker://ghcr.io/nevotheless/omar:rolling", "ghcr.io/nevotheless/omar:rolling"},
		{"ubuntu", "ubuntu:latest"},
	}
	for _, tt := range tests {
		got := normalizeImageRef(tt.in)
		if got != tt.out {
			t.Errorf("normalizeImageRef(%q) = %q, want %q", tt.in, got, tt.out)
		}
	}
}
