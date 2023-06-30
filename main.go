package main

import (
	"github.com/spf13/cobra"
	"github.com/zbysir/writeflow/cmd"
	"github.com/zbysir/writeflow/internal/pkg/log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:           "writeflow",
	Short:         "writeflow",
	Long:          `writeflow`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(cmd.Api())
	rootCmd.AddCommand(cmd.Tool())
	rootCmd.AddCommand(cmd.Version("v0.0.1"))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Errorf("%+v", err)
		os.Exit(1)
	}
}
