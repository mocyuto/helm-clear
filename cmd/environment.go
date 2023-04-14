package main

import (
	"github.com/spf13/pflag"
)

type EnvSettings struct {
	DryRun         bool
	KubeConfigFile string
	KubeContext    string
	KubeNamespace  string
}

func (s *EnvSettings) AddBaseFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.DryRun, "dry-run", false, "simulate an install")
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	s.AddBaseFlags(fs)
	fs.StringVar(&s.KubeConfigFile, "kubeconfig", "", "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVar(&s.KubeNamespace, "namespace", s.KubeNamespace, "namespace scope for this request")
}
