package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nevotheless/omar/internal/image"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show CLI and image version",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := image.Current()
			if err != nil {
				fmt.Printf("omar CLI version: %s\n", image.CLIVersion())
				return nil
			}
			fmt.Printf("omar CLI version:   %s\n", image.CLIVersion())
			fmt.Printf("System image:       %s\n", info.Image)
			if info.Version != "" {
				fmt.Printf("Image version:      %s\n", info.Version)
			}
			return nil
		},
	}
}
