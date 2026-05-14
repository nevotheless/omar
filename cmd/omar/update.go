package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check for and apply image updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Checking for updates...")
			return nil
		},
	}
}
