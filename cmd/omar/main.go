package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "omar",
		Short: "Immutable OCI-image-based OS – atomic updates for omarchy",
		Long: `omar manages atomic, image-based updates for omarchy.

It converts existing mutable systems to immutable, handles daily OCI image
updates via bootc, and wraps Flatpak/Distrobox for user package management.`,
		SilenceUsage: true,
	}

	root.AddCommand(newInstallCmd())
	root.AddCommand(newUpdateCmd())
	root.AddCommand(newRollbackCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newPkgCmd())
	root.AddCommand(newVersionCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
