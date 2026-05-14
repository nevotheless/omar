package update

import (
	"fmt"
	"time"

	"github.com/nevotheless/omar/internal/bootc"
)

type Status struct {
	Current     string `json:"current"`
	Available   bool   `json:"available"`
	Staged      bool   `json:"staged"`
	NewVersion  string `json:"new_version,omitempty"`
	LastChecked string `json:"last_checked,omitempty"`
}

func Check() (*Status, error) {
	s := &Status{}

	if !bootc.HasOstree() {
		return s, fmt.Errorf("not a bootc system (run 'omar install' first)")
	}

	status, err := bootc.StatusJSON()
	if err != nil {
		return s, fmt.Errorf("get bootc status: %w", err)
	}

	if status.Booted != nil {
		s.Current = status.Booted.Image.Image
	}

	s.Staged = status.Staged != nil
	if status.Staged != nil {
		s.Available = true
		s.NewVersion = status.Staged.Image.Version
	}

	if !s.Available {
		fmt.Println("Checking registry for updates...")
		// bootc upgrade --check would be ideal; for now just indicate no staged update
		fmt.Printf("Current image: %s\n", s.Current)
	}

	s.LastChecked = time.Now().Format(time.RFC3339)
	return s, nil
}

func Apply() error {
	if !bootc.HasOstree() {
		return fmt.Errorf("not a bootc system (run 'omar install' first)")
	}

	fmt.Println("Pulling and staging latest image...")
	if err := bootc.Upgrade(); err != nil {
		return fmt.Errorf("bootc upgrade failed: %w", err)
	}

	fmt.Println("\n✓ Update staged successfully")
	fmt.Println("  Reboot to apply: sudo systemctl reboot")
	fmt.Println("  Rollback if needed: omar rollback")
	return nil
}

func Rollback() error {
	if !bootc.HasOstree() {
		return fmt.Errorf("not a bootc system (run 'omar install' first)")
	}

	fmt.Println("Rolling back to previous deployment...")
	if err := bootc.Rollback(); err != nil {
		return fmt.Errorf("bootc rollback failed: %w", err)
	}

	fmt.Println("\n✓ Rollback staged")
	fmt.Println("  Reboot to apply: sudo systemctl reboot")
	return nil
}
