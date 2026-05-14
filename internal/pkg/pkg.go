package pkg

import (
	"fmt"
	"os/exec"
	"strings"
)

type Backend int

const (
	BackendUnknown Backend = iota
	BackendFlatpak
	BackendDistrobox
)

var guiPackages = map[string]bool{
	"firefox":          true,
	"chromium":         true,
	"brave":            true,
	"spotify":          true,
	"discord":          true,
	"slack":            true,
	"telegram-desktop": true,
	"code":             true,
	"obsidian":         true,
	"gimp":             true,
	"inkscape":         true,
	"vlc":              true,
	"celluloid":        true,
	"evince":           true,
	"thunderbird":      true,
	"onlyoffice":       true,
	"libreoffice":      true,
	"zeal":             true,
	"obs-studio":       true,
}

func Detect(pkg string) Backend {
	if guiPackages[pkg] {
		return BackendFlatpak
	}

	flatpakOk := hasFlatpak()
	distroboxOk := hasDistrobox()

	if flatpakOk && !distroboxOk {
		return BackendFlatpak
	}
	if distroboxOk && !flatpakOk {
		return BackendDistrobox
	}

	return BackendDistrobox
}

func Install(pkg string) error {
	backend := Detect(pkg)
	fmt.Printf("Installing %s via %s...\n", pkg, backendLabel(backend))

	switch backend {
	case BackendFlatpak:
		return installFlatpak(pkg)
	case BackendDistrobox:
		return installDistrobox(pkg)
	default:
		return fmt.Errorf("no suitable backend for %s", pkg)
	}
}

func List() ([]string, error) {
	var result []string

	flatpaks, _ := listFlatpak()
	result = append(result, flatpaks...)

	distroboxPkgs, _ := listDistrobox()
	result = append(result, distroboxPkgs...)

	return result, nil
}

func hasFlatpak() bool {
	_, err := exec.LookPath("flatpak")
	return err == nil
}

func hasDistrobox() bool {
	_, err := exec.LookPath("distrobox")
	return err == nil
}

func installFlatpak(pkg string) error {
	cmd := exec.Command("flatpak", "install", "--user", "-y", "flathub", pkg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func installDistrobox(pkg string) error {
	containerName := "dev"
	_ = ensureDistroboxContainer(containerName)
	cmd := exec.Command("distrobox", "enter", containerName, "--", "sudo", "pacman", "-S", "--noconfirm", pkg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func ensureDistroboxContainer(name string) error {
	cmd := exec.Command("distrobox", "list")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if strings.Contains(string(out), name) {
		return nil
	}
	fmt.Printf("Creating distrobox container '%s'...\n", name)
	create := exec.Command("distrobox", "create", "--name", name, "--image", "archlinux:latest")
	create.Stdout = nil
	create.Stderr = nil
	return create.Run()
}

func listFlatpak() ([]string, error) {
	cmd := exec.Command("flatpak", "list", "--app", "--columns=application")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var pkgs []string
	for _, l := range lines {
		if l != "" && l != "Application ID" {
			pkgs = append(pkgs, "flatpak:"+l)
		}
	}
	return pkgs, nil
}

func listDistrobox() ([]string, error) {
	cmd := exec.Command("distrobox", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var pkgs []string
	for _, l := range lines {
		if l != "" && !strings.HasPrefix(l, "ID") {
			pkgs = append(pkgs, "distrobox:"+strings.Fields(l)[0])
		}
	}
	return pkgs, nil
}

func backendLabel(b Backend) string {
	switch b {
	case BackendFlatpak:
		return "Flatpak"
	case BackendDistrobox:
		return "Distrobox"
	default:
		return "unknown"
	}
}
