package main

import (
	"fmt"

	"github.com/nevotheless/omar/internal/update"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check for and apply image updates",
		Long: `Pulls the latest OCI image from the registry, stages it atomically,
and prompts for reboot. The previous image remains as rollback option.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := update.Check()
			if err != nil {
				return err
			}

			fmt.Printf("Current: %s\n", status.Current)
			_ = status

			fmt.Println()
			if err := update.Apply(); err != nil {
				return err
			}
			return nil
		},
	}
}
