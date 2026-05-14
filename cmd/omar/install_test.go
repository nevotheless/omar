package main

import (
	"testing"
)

func TestFreshInstall_MissingDisk(t *testing.T) {
	err := freshInstall("", "test-image", false)
	if err == nil {
		t.Fatal("expected error for empty disk")
	}
}

func TestFreshInstall_NoBootc(t *testing.T) {
	oldCheck := checkBootc
	checkBootc = func() bool { return false }
	defer func() { checkBootc = oldCheck }()

	oldConfirm := confirmDestructive
	confirmDestructive = func(disk string) error { return nil }
	defer func() { confirmDestructive = oldConfirm }()

	err := freshInstall("/dev/sda", "test-image", false)
	if err == nil {
		t.Fatal("expected error when bootc is not installed")
	}
}

func TestFreshInstall_AutoYesSkipsConfirm(t *testing.T) {
	oldCheck := checkBootc
	checkBootc = func() bool { return true }
	defer func() { checkBootc = oldCheck }()

	oldConfirm := confirmDestructive
	called := false
	confirmDestructive = func(disk string) error {
		called = true
		return nil
	}
	defer func() { confirmDestructive = oldConfirm }()

	err := freshInstall("/dev/sda", "test-image", true)
	if err == nil {
		t.Fatal("expected error from bootc.Install (not installed in test env)")
	}
	if called {
		t.Fatal("confirmDestructive should not be called with autoYes=true")
	}
}

func TestFreshInstall_RequiresConfirm(t *testing.T) {
	oldCheck := checkBootc
	checkBootc = func() bool { return true }
	defer func() { checkBootc = oldCheck }()

	oldConfirm := confirmDestructive
	called := false
	confirmDestructive = func(disk string) error {
		called = true
		return nil
	}
	defer func() { confirmDestructive = oldConfirm }()

	err := freshInstall("/dev/sda", "test-image", false)
	if err == nil {
		t.Fatal("expected error from bootc.Install (not installed in test env)")
	}
	if !called {
		t.Fatal("confirmDestructive should be called with autoYes=false")
	}
}
