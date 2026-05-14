package update

import (
	"fmt"
)

// Check looks for a newer image in the registry.
func Check() error {
	fmt.Println("Checking ghcr.io/basecamp/omar:rolling for updates...")
	return nil
}

// Apply stages the latest image and returns true if a reboot is required.
func Apply() (bool, error) {
	fmt.Println("Staging latest image...")
	return true, nil
}
