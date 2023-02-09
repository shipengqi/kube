package kube

import (
	"flag"
	"os"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	kubeconfig    string
	podname       string
	containername string
	namespace     string
	mockcli       *Client
)

type testcm struct {
	name string
	data map[string]string
}

func TestMain(m *testing.M) {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file")
	flag.StringVar(&podname, "pod", "", "pod name for exec")
	flag.StringVar(&containername, "container", "", "container name for exec")
	flag.StringVar(&namespace, "namespace", "", "namespace for exec")
	flag.Parse()

	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &kubeconfig
	cfg := NewConfig(flags)
	mockcli = New(cfg)

	os.Exit(m.Run())
}
