package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var image string
	var fresh bool
	var disk string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Convert mutable system to immutable, or fresh-install to disk",
		Long: `Without flags: converts the running mutable omarchy to an immutable
bootc-based system using 'bootc switch'.

With --fresh: installs omar to an empty disk using 'bootc install'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if fresh {
				if disk == "" {
					return fmt.Errorf("--fresh requires a target disk (e.g. /dev/sda)")
				}
				fmt.Printf("Fresh install to %s from %s\n", disk, image)
				return nil
			}
			fmt.Printf("Converting mutable system to immutable using %s\n", image)
			return nil
		},
	}

	cmd.Flags().StringVar(&image, "from", "ghcr.io/basecamp/omar:rolling", "OCI image reference")
	cmd.Flags().BoolVar(&fresh, "fresh", false, "Fresh install to disk")
	cmd.Flags().StringVar(&disk, "disk", "", "Target disk for fresh install")
	return cmd
}
