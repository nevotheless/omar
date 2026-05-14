package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "omar",
		Short: "Immutable OCI-image-based omarchy OS",
		Long: `omar manages atomic, image-based updates for omarchy.
It converts existing mutable systems to immutable, handles daily OCI image updates,
and wraps Flatpak/Distrobox for user package management.`,
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
