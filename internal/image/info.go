package image

import (
	"fmt"

	"github.com/nevotheless/omar/internal/bootc"
	"github.com/nevotheless/omar/internal/version"
)

type Info struct {
	Image          string `json:"image"`
	Version        string `json:"version"`
	DeploymentID   string `json:"deployment_id"`
	Booted         bool   `json:"booted"`
	Staged         bool   `json:"staged"`
	RollbackExists bool   `json:"rollback_exists"`
}

func Current() (*Info, error) {
	info := &Info{}

	if !bootc.HasOstree() {
		info.Image = "none (mutable system)"
		return info, nil
	}

	status, err := bootc.StatusJSON()
	if err != nil {
		return info, fmt.Errorf("get bootc status: %w", err)
	}

	if status.Booted != nil {
		info.Image = status.Booted.Image.Image
		info.Version = status.Booted.Image.Version
		info.DeploymentID = status.Booted.ID
		info.Booted = true
	}
	info.Staged = status.Staged != nil
	info.RollbackExists = status.Rollback != nil

	return info, nil
}

func CLIVersion() string {
	return version.Version
}

func String() string {
	i, err := Current()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	s := fmt.Sprintf("Image:       %s\n", i.Image)
	s += fmt.Sprintf("Version:     %s\n", i.Version)
	if i.Booted {
		s += fmt.Sprintf("Deployment:  %s\n", i.DeploymentID)
	}
	s += fmt.Sprintf("Booted:      %v\n", i.Booted)
	s += fmt.Sprintf("Staged:      %v\n", i.Staged)
	s += fmt.Sprintf("Rollback:    %v\n", i.RollbackExists)
	return s
}
