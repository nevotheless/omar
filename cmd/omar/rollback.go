package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nevotheless/omar/internal/update"
)

func newRollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback to previous deployment",
		Long: `Stages the previous deployment for boot. After reboot the system
will be running the prior image.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := update.Rollback(); err != nil {
				return err
			}
			fmt.Println("Run 'sudo systemctl reboot' to apply the rollback.")
			return nil
		},
	}
}
