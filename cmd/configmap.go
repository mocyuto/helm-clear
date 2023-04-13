package main

import (
	"errors"
	"io"
	"os"

	"github.com/mocyuto/helm-clear/pkg"
	"github.com/spf13/cobra"
)

var settings *EnvSettings

func newConfigmapCmd(_ io.Writer, args []string) *cobra.Command {
	options := &pkg.ConfigmapOptions{}
	cmd := &cobra.Command{
		Use:   "configmap [flags] ChartName",
		Short: "clear configmaps",
		Long:  `clear configmaps`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(1)
			}
			if len(args) != 2 {
				return errors.New("wrong number of arguments")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options.DryRun = settings.DryRun
			options.Namespace = settings.KubeNamespace
			options.ChartName = args[1]
			if err := pkg.RunConfigmap(*options); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	_ = flags.Parse(args)

	settings = new(EnvSettings)

	if ctx := os.Getenv("HELM_KUBECONTEXT"); ctx != "" {
		settings.KubeContext = ctx
	}
	if ns := os.Getenv("HELM_NAMESPACE"); ns != "" {
		settings.KubeNamespace = ns
	}

	// Note that the plugin's --kubeconfig flag is set by the Helm plugin framework to
	// the KUBECONFIG environment variable instead of being passed into the plugin.
	settings.AddFlags(flags)

	cmd.Flags().IntVar(&options.History, "history", 5, "retain history count")
	cmd.Flags().StringVar(&options.ChartName, "name", "", "chart name")
	cmd.Flags().StringVar(&options.ConfigmapName, "configmap", "", "configmap names")

	return cmd
}
