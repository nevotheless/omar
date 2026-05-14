package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback to previous deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Rolling back to previous deployment...")
			return nil
		},
	}
}
