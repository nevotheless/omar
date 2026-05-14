package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/nevotheless/omar/internal/bootc"
)

type Plan struct {
	Image      string
	HasESP     bool
	HasBootctl bool
	HasBootc   bool
	FreeGB     int64
	CanConvert bool
	Issues     []string
}

const DefaultImage = "ghcr.io/nevotheless/omar:rolling"
const minFreeGB = 20

func Check(image string) (*Plan, error) {
	if image == "" {
		image = DefaultImage
	}
	p := &Plan{Image: image}

	if bootc.HasOstree() {
		p.Issues = append(p.Issues, "system is already running bootc/ostree")
		return p, nil
	}

	p.HasESP = checkESP()
	p.HasBootctl = checkBootctl()
	p.HasBootc = bootc.HasBootc()
	p.FreeGB = freeSpaceGB()

	if !p.HasESP {
		p.Issues = append(p.Issues, "no EFI system partition found (mount point /boot or /efi with vfat)")
	}
	if !p.HasBootctl {
		p.Issues = append(p.Issues, "systemd-boot not installed (run: bootctl install)")
	}
	if p.FreeGB < minFreeGB {
		p.Issues = append(p.Issues, fmt.Sprintf("only %dGB free on root, need at least %dGB", p.FreeGB, minFreeGB))
	}

	p.CanConvert = len(p.Issues) == 0
	return p, nil
}

func Migrate(image string, autoYes bool) error {
	plan, err := Check(image)
	if err != nil {
		return err
	}

	fmt.Println("=== Pre-flight Check ===")
	fmt.Printf("  Image:         %s\n", plan.Image)
	fmt.Printf("  ESP:           %v\n", plan.HasESP)
	fmt.Printf("  systemd-boot:  %v\n", plan.HasBootctl)
	fmt.Printf("  bootc:         %v\n", plan.HasBootc)
	fmt.Printf("  Free space:    %d GB\n", plan.FreeGB)

	if !plan.CanConvert {
		fmt.Println("\n✗ Pre-flight checks failed:")
		for _, issue := range plan.Issues {
			fmt.Printf("  - %s\n", issue)
		}
		return fmt.Errorf("conversion pre-flight checks failed")
	}

	if !plan.HasBootc {
		fmt.Println("\nInstalling bootc...")
		if err := bootc.InstallBootc(); err != nil {
			return fmt.Errorf("install bootc: %w", err)
		}
	}

	fmt.Println("\n=== Starting Conversion ===")
	fmt.Printf("Switching to immutable image: %s\n", plan.Image)

	if err := bootc.Switch(plan.Image); err != nil {
		return fmt.Errorf("bootc switch failed: %w", err)
	}

	fmt.Println("\n✓ System prepared for immutability")
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println("    1. Run 'omar update' to pull latest image")
	fmt.Println("    2. Reboot: sudo systemctl reboot")
	fmt.Println("    3. After reboot, verify with 'omar status'")
	fmt.Println()
	fmt.Println("  The previous system remains as rollback entry in systemd-boot.")
	return nil
}

func checkESP() bool {
	for _, mount := range []string{"/boot", "/efi"} {
		if isVfat(mount) {
			return true
		}
	}
	return false
}

func isVfat(path string) bool {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return false
	}
	// VFAT super magic: 0x4d44
	return stat.Type == 0x4d44
}

func checkBootctl() bool {
	if _, err := os.Stat("/usr/lib/systemd/boot/efi/systemd-bootx64.efi"); err == nil {
		return true
	}
	if _, err := os.Stat("/usr/lib/systemd-bootx64.efi"); err == nil {
		return true
	}
	return false
}

func freeSpaceGB() int64 {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return 0
	}
	free := int64(stat.Bavail) * int64(stat.Bsize)
	return free / (1 << 30)
}

func IsSystemdBooted() (bool, error) {
	body, err := os.ReadFile("/proc/1/comm")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(body)) == "systemd", nil
}

func findESP() string {
	for _, path := range []string{"/boot", "/efi", "/boot/efi"} {
		info, err := os.Stat(filepath.Join(path, "EFI"))
		if err == nil && info.IsDir() {
			return path
		}
	}
	return ""
}
