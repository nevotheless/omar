package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nevotheless/omar/internal/image"
)

type versionOutput struct {
	CliVersion string `json:"cli_version"`
	Image      string `json:"image,omitempty"`
	Version    string `json:"version,omitempty"`
	Booted     bool   `json:"booted"`
	Error      string `json:"error,omitempty"`
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show CLI and image version",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonFlag, _ := cmd.Flags().GetBool("json")

			info, err := image.Current()

			if jsonFlag {
				out := versionOutput{CliVersion: image.CLIVersion()}
				if err == nil {
					out.Image = info.Image
					out.Version = info.Version
					out.Booted = info.Booted
				} else {
					out.Error = err.Error()
				}
				b, marshalErr := json.MarshalIndent(out, "", "  ")
				if marshalErr != nil {
					return marshalErr
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				return nil
			}

			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "omar CLI version: %s\n", image.CLIVersion())
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "omar CLI version:   %s\n", image.CLIVersion())
			fmt.Fprintf(cmd.OutOrStdout(), "System image:       %s\n", info.Image)
			if info.Version != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Image version:      %s\n", info.Version)
			}
			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}
