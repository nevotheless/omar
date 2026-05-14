package pkg

import "fmt"

// Backend represents a package installation backend.
type Backend int

const (
	BackendUnknown Backend = iota
	BackendFlatpak
	BackendDistrobox
)

// InstallKind detects whether a package is GUI or CLI and installs accordingly.
func Install(name string) error {
	backend, err := detect(name)
	if err != nil {
		return err
	}
	switch backend {
	case BackendFlatpak:
		return installFlatpak(name)
	case BackendDistrobox:
		return installDistrobox(name)
	default:
		return fmt.Errorf("unknown backend for %s", name)
	}
}

func detect(pkg string) (Backend, error) {
	// TODO: heuristic – check known GUI vs CLI packages, or try flatpak first
	return BackendDistrobox, nil
}

func installFlatpak(pkg string) error {
	fmt.Printf("flatpak install flathub %s\n", pkg)
	return nil
}

func installDistrobox(pkg string) error {
	fmt.Printf("distrobox enter --additional-flags \"sudo pacman -S %s\"\n", pkg)
	return nil
}
