# kube

A simple Kubernetes client, based on [client-go](https://github.com/kubernetes/client-go).

## Quick Start

```go
package main

import (
    "context"
    "log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	
	"github.com/shipengqi/kube"
	
)

func main() {
	kubeconfig := "testdata/config"
	
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &kubeconfig
	cfg := kube.NewConfig(flags)
	cli := kube.New(cfg)
	k8s, err := cli.Dial()
	if err != nil {
		log.Fatal(err)
	}
	
	// get a configmap named "configmapname"
	cm, err := k8s.CoreV1().ConfigMaps("default").Get(context.TODO(), "configmapname", metav1.GetOptions{})
	log.Println(cm.Data)
	
	// or 
	cm, err = cli.GetConfigMap(context.TODO(), "default", "configmapname")
	log.Println(cm.Data)
	
	// apply file, is like "kubectl apply -f testdata/content-apply.yaml"
	err = cli.Apply([]string{"testdata/content-apply.yaml"})
	if err != nil {
		log.Fatal(err)
	}

	// Exec in a pod, is like "kubectl exec <pod name> -n <namespace> -c <container name> -- <command>"
	stdout, stderr, err := cli.Exec("podname", "containername", "namespace", "command")
	if err != nil {
		log.Println(stderr)
		log.Fatal(err)
	}
	log.Println(stdout)
}
```

## Test

```bash
go test -v . -kubeconfig <kubeconfig file>
```

### Test Client.Exec
```bash
go test -v -kubeconfig <kubeconfig file> -container <container name> -pod <pod name> -namespace <namespace> .
```
