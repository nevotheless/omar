package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/nevotheless/omar/internal/bootc"
	"github.com/nevotheless/omar/internal/convert"
)

func newInstallCmd() *cobra.Command {
	var image string
	var fresh bool
	var disk string
	var autoYes bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Convert mutable system to immutable, or fresh-install to disk",
		Long: `Without flags: converts the running mutable omarchy to an immutable
bootc-based system using 'bootc switch'. Pre-flight checks verify that
the system has an ESP, systemd-boot, and enough free space.

With --fresh: installs omar to an empty disk using 'bootc install'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if fresh {
				if disk == "" {
					return fmt.Errorf("--fresh requires a target disk (e.g. /dev/sda)")
				}
			return freshInstall(disk, image, autoYes)
		}
		return convert.Migrate(image, autoYes)
		},
	}

	cmd.Flags().StringVar(&image, "from", convert.DefaultImage, "OCI image reference")
	cmd.Flags().BoolVar(&fresh, "fresh", false, "Fresh install to disk (requires --disk)")
	cmd.Flags().StringVar(&disk, "disk", "", "Target disk for fresh install")
	cmd.Flags().BoolVarP(&autoYes, "yes", "y", false, "Skip confirmation prompts")
	return cmd
}

func freshInstall(disk, image string, autoYes bool) error {
	if disk == "" {
		return fmt.Errorf("--fresh requires a target disk (e.g. /dev/sda)")
	}
	if !checkBootc() {
		return fmt.Errorf("bootc is not installed. Run 'omar install' without --fresh first.")
	}

	fmt.Printf("This will DESTROY all data on %s\n", disk)
	fmt.Printf("Installing from image: %s\n", image)
	fmt.Println()

	if !autoYes {
		if err := confirmDestructive(disk); err != nil {
			return err
		}
	}

	return bootc.Install(disk, image)
}

var checkBootc = bootc.HasBootc

var confirmDestructive = defaultConfirm

func defaultConfirm(disk string) error {
	fmt.Printf("Type the disk name (%s) to confirm: ", disk)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return fmt.Errorf("aborted")
	}
	if input != disk {
		return fmt.Errorf("confirmation mismatch, aborted")
	}
	return nil
}

func init() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		confirmDestructive = func(disk string) error { return nil }
	}
}
