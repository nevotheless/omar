package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current deployment and image info",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("System status:")
			fmt.Println("  Image:       <not deployed>")
			fmt.Println("  Deployments: 0")
			fmt.Println("  Next boot:   <current>")
			return nil
		},
	}
}
