package update

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/nevotheless/omar/internal/bootc"
)

// execCommand allows mocking in tests
var execCommand = exec.Command

// execLookPath allows mocking in tests
var execLookPath = exec.LookPath

type Status struct {
	Current     string `json:"current"`
	Available   bool   `json:"available"`
	Staged      bool   `json:"staged"`
	NewVersion  string `json:"new_version,omitempty"`
	LastChecked string `json:"last_checked,omitempty"`
	Registry    string `json:"registry,omitempty"`
}

// bootcUpgradeCheckResult holds the parsed result of `bootc upgrade --check`
type bootcUpgradeCheckResult struct {
	UpdateAvailable bool   `json:"update_available"`
	Image           string `json:"image,omitempty"`
	Version         string `json:"version,omitempty"`
	Digest          string `json:"digest,omitempty"`
	Message         string `json:"message,omitempty"`
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
		s.Registry = status.Booted.Image.Image
	}

	s.Staged = status.Staged != nil
	if status.Staged != nil {
		s.Available = true
		s.NewVersion = status.Staged.Image.Version
	}

	if !s.Available && s.Current != "" {
		fmt.Println("Checking registry for updates...")
		newVersion, available, err := checkRegistry(s.Current)
		if err != nil {
			fmt.Printf("  Registry check note: %v\n", err)
		}
		if available {
			s.Available = true
			s.NewVersion = newVersion
		}
		fmt.Printf("Current image: %s\n", s.Current)
	}

	s.LastChecked = time.Now().Format(time.RFC3339)
	return s, nil
}

// checkRegistry attempts to detect a newer image version from the registry.
// Strategy:
//  1. Try `bootc upgrade --check` (primary, bootc v1.1+)
//  2. Fallback: try `skopeo inspect` to read manifest labels
func checkRegistry(currentImage string) (string, bool, error) {
	// Strategy 1: bootc upgrade --check
	if newVer, ok, err := tryBootcUpgradeCheck(); err == nil && ok {
		return newVer, true, nil
	} else if err == nil {
		// bootc says no update — trusted
		return "", false, nil
	}
	// bootc upgrade --check failed or unavailable, fall through

	// Strategy 2: skopeo inspect
	if hasSkopeo() {
		return trySkopeoCheck(currentImage)
	}

	return "", false, fmt.Errorf("no update check method available (bootc upgrade --check and skopeo unavailable)")
}

func tryBootcUpgradeCheck() (string, bool, error) {
	cmd := execCommand("bootc", "upgrade", "--check")
	out, err := cmd.Output()
	if err != nil {
		// bootc doesn't support --check or other error
		return "", false, fmt.Errorf("bootc upgrade --check failed: %w", err)
	}

	output := strings.TrimSpace(string(out))

	// Try parsing as JSON (bootc v1.1+ JSON output)
	var result bootcUpgradeCheckResult
	if err := json.Unmarshal([]byte(output), &result); err == nil {
		if result.UpdateAvailable {
			return result.Version, true, nil
		}
		return "", false, nil
	}

	// Fallback: plain text parsing
	if strings.Contains(output, "No update available") ||
		strings.Contains(output, "already up to date") ||
		strings.Contains(output, "already deployed") {
		return "", false, nil
	}

	// Any other non-empty output likely means an update is available
	if output != "" {
		return output, true, nil
	}

	return "", false, nil
}

func trySkopeoCheck(currentImage string) (string, bool, error) {
	// extract registry/repo:tag from image reference
	img := normalizeImageRef(currentImage)

	cmd := execCommand("skopeo", "inspect", "--raw", "docker://"+img)
	out, err := cmd.Output()
	if err != nil {
		return "", false, fmt.Errorf("skopeo inspect failed: %w", err)
	}

	// Parse OCI manifest
	var manifest struct {
		Config struct {
			Labels map[string]string `json:"labels"`
		} `json:"config"`
		Annotations map[string]string `json:"annotations"`
	}
	if err := json.Unmarshal(out, &manifest); err != nil {
		return "", false, fmt.Errorf("parse manifest: %w", err)
	}

	// Check version labels in order of preference
	version := ""
	if manifest.Config.Labels != nil {
		version = manifest.Config.Labels["version"]
		if version == "" {
			version = manifest.Config.Labels["org.opencontainers.image.version"]
		}
	}
	if version == "" && manifest.Annotations != nil {
		version = manifest.Annotations["org.opencontainers.image.version"]
	}

	if version != "" {
		return version, true, nil
	}

	return "", false, fmt.Errorf("no version label found in manifest")
}

func hasSkopeo() bool {
	_, err := execLookPath("skopeo")
	return err == nil
}

// normalizeImageRef ensures the image reference is suitable for skopeo.
// e.g. "ghcr.io/nevotheless/omar:rolling" stays as is.
func normalizeImageRef(ref string) string {
	// Remove leading docker:// if present
	ref = strings.TrimPrefix(ref, "docker://")
	// Remove tag if not present, default to :latest
	if !strings.Contains(ref, ":") {
		ref += ":latest"
	}
	return ref
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
