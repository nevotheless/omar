package main

import (
	"fmt"

	"github.com/nevotheless/omar/internal/pkg"
	"github.com/spf13/cobra"
)

func newPkgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pkg",
		Short: "Manage user packages (Flatpak, Distrobox)",
		Long: `Smart package installation: GUI apps go to Flatpak,
CLI/dev tools go to a Distrobox container.
System packages must be baked into the OCI image.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add <package>",
		Short: "Install a package (GUI→Flatpak, CLI→Distrobox)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pkg.Install(args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List installed packages and their origin",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkgs, err := pkg.List()
			if err != nil {
				return err
			}
			if len(pkgs) == 0 {
				fmt.Println("No user packages installed.")
				return nil
			}
			fmt.Println("Installed packages:")
			for _, p := range pkgs {
				fmt.Printf("  %s\n", p)
			}
			return nil
		},
	})

	return cmd
}
