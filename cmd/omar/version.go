package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show CLI and image version",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliVer := "0.1.0-dev"
			fmt.Printf("omar CLI version: %s\n", cliVer)
			return nil
		},
	}
}
