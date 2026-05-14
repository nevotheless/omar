package bootc

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Status struct {
	Booted   *Deployment  `json:"booted"`
	Staged   *Deployment  `json:"staged,omitempty"`
	Rollback *Deployment  `json:"rollback,omitempty"`
}

type Deployment struct {
	ID        string   `json:"id"`
	Image     ImageRef `json:"image"`
	Timestamp int64    `json:"timestamp"`
}

type ImageRef struct {
	Image   string `json:"image"`
	Version string `json:"version"`
}

func StatusJSON() (*Status, error) {
	cmd := exec.Command("bootc", "status", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("bootc status: %w", err)
	}
	var s Status
	if err := json.Unmarshal(out, &s); err != nil {
		return nil, fmt.Errorf("parse bootc status: %w", err)
	}
	return &s, nil
}

func Switch(image string) error {
	cmd := exec.Command("bootc", "switch", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Upgrade() error {
	cmd := exec.Command("bootc", "upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Rollback() error {
	cmd := exec.Command("bootc", "rollback")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Install(disk, image string) error {
	args := []string{"install", "--target-disk", disk}
	if image != "" {
		args = append(args, image)
	}
	cmd := exec.Command("bootc", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func HasBootc() bool {
	_, err := exec.LookPath("bootc")
	return err == nil
}

func HasOstree() bool {
	_, err := os.Stat("/run/ostree-booted")
	return err == nil
}

func InstallBootc() error {
	fmt.Println("Installing bootc from AUR...")
	cmd := exec.Command("yay", "-S", "--noconfirm", "bootc-bin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
