package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Version(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s\n", version)
			return nil
		},
	}
}