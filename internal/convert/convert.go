package convert

import (
	"errors"
	"fmt"
)

// Plan holds the pre-flight check results for a mutable→immutable conversion.
type Plan struct {
	Image      string
	HasESP     bool
	HasBootctl bool
	HasBootc   bool
	FreeGB     int
	CanConvert bool
	Issues     []string
}

// Check runs pre-flight checks on the current system.
func Check(image string) (*Plan, error) {
	p := &Plan{Image: image}

	// TODO: implement actual checks
	p.HasESP = true
	p.HasBootctl = true
	p.HasBootc = false
	p.FreeGB = 50
	p.CanConvert = true

	if !p.HasESP {
		p.Issues = append(p.Issues, "no EFI system partition found")
	}
	if !p.HasBootctl {
		p.Issues = append(p.Issues, "systemd-boot not installed")
	}
	if p.FreeGB < 20 {
		p.Issues = append(p.Issues, fmt.Sprintf("only %dGB free, need at least 20GB", p.FreeGB))
	}
	if len(p.Issues) > 0 {
		p.CanConvert = false
	}

	return p, nil
}

// Migrate converts the running system to an immutable bootc deployment.
func Migrate(image string) error {
	plan, err := Check(image)
	if err != nil {
		return err
	}
	if !plan.CanConvert {
		return errors.New("conversion pre-flight checks failed")
	}
	// TODO: install bootc if missing, then run bootc switch
	return nil
}
