package main

import "os"

func main() {
	configmapCmd := newConfigmapCmd(os.Stdout, os.Args[1:])
	if err := configmapCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
