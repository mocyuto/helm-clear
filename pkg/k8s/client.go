package k8s

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientset "k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewClient(kubeCtx string) (*clientset.Clientset, error) {
	getter := genericclioptions.NewConfigFlags(true)
	if kubeCtx != "" {
		getter.Context = &kubeCtx
	}
	factory := cmdutil.NewFactory(getter)
	k, err := factory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	return k, nil
}
