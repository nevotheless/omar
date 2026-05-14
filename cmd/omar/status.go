package main

import (
	"fmt"

	"github.com/nevotheless/omar/internal/image"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current deployment and image info",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := image.Current()
			if err != nil {
				return err
			}

			fmt.Println("=== omar System Status ===")
			fmt.Printf("  Image:       %s\n", info.Image)
			fmt.Printf("  Version:     %s\n", info.Version)
			if info.Booted {
				fmt.Printf("  Deployment:  %s\n", info.DeploymentID)
			}
			fmt.Printf("  Booted:      %v\n", info.Booted)
			fmt.Printf("  Staged:      %v\n", info.Staged)
			fmt.Printf("  Rollback:    %v\n", info.RollbackExists)
			fmt.Printf("  CLI version: %s\n", image.CLIVersion())
			return nil
		},
	}
}
