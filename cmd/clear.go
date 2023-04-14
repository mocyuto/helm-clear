package main

import (
	"fmt"
	"io"
	"os"

	helm_clear "github.com/mocyuto/helm-clear"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clear",
	Short: "helm plugin for clear old revision resources",
}

func main() {
	rootCmd.AddCommand(
		newConfigmapCmd(os.Stdout, os.Args[1:]),
		newVersionCmd(os.Stdout),
	)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newVersionCmd(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print the version number of helm-clear",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(out, helm_clear.Version)
		},
	}
}
