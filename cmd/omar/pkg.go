package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPkgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pkg",
		Short: "Manage user packages (Flatpak, Distrobox)",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add <package>",
		Short: "Install a package (GUI→Flatpak, CLI→Distrobox)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("package name required")
			}
			fmt.Printf("Installing %s...\n", args[0])
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List installed packages and their origin",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Installed packages:")
			return nil
		},
	})

	return cmd
}
