package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	cmd := newVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("version cmd without --json failed: %v", err)
	}
}

func TestVersionCmdJSON(t *testing.T) {
	cmd := newVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("version cmd with --json failed: %v", err)
	}
	if !json.Valid(buf.Bytes()) {
		t.Fatalf("output is not valid JSON:\n%s", buf.String())
	}
	var out versionOutput
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v\nraw: %s", err, buf.String())
	}
	if out.CliVersion == "" {
		t.Error("expected cli_version to be non-empty")
	}
}
