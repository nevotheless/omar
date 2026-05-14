package image

import (
	"fmt"
	"os"
)

// Info holds details about the currently deployed image.
type Info struct {
	Image      string
	Version    string
	Deployment string
}

// Current returns the status of the running deployment.
func Current() (*Info, error) {
	// bootc status --json would be the real source
	if _, err := os.Stat("/run/ostree-booted"); os.IsNotExist(err) {
		return &Info{Image: "none (mutable system)"}, nil
	}
	return &Info{
		Image:      "ghcr.io/basecamp/omar:rolling",
		Version:    "rolling-20260514",
		Deployment: "0",
	}, nil
}

func String() string {
	info, err := Current()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("Image: %s\nVersion: %s\nDeployment: %s",
		info.Image, info.Version, info.Deployment)
}
