# kube

A simple Kubernetes client, based on [client-go](https://github.com/kubernetes/client-go).

[![E2E](https://github.com/shipengqi/kube/actions/workflows/e2e.yaml/badge.svg)](https://github.com/shipengqi/kube/actions/workflows/e2e.yaml)
[![codecov](https://codecov.io/gh/shipengqi/kube/branch/main/graph/badge.svg?token=0KSRZKV4C8)](https://codecov.io/gh/shipengqi/kube)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipengqi/kube)](https://goreportcard.com/report/github.com/shipengqi/kube)
[![Release](https://img.shields.io/github/release/shipengqi/kube.svg)](https://github.com/shipengqi/kube/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipengqi/kube)](https://goreportcard.com/report/github.com/shipengqi/kube)
[![License](https://img.shields.io/github/license/shipengqi/kube)](https://github.com/shipengqi/kube/blob/main/LICENSE)

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
	err = cli.Apply("testdata/content-apply.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// delete file, is like "kubectl delete -f testdata/content-apply.yaml"
	err = cli.Delete("testdata/content-apply.yaml")
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

	// Upload uploads local file to a remote pod
	err = cli.Upload(context.TODO(), "podname", "containername", "namespace", "testdata/upload.txt", "/testdata")
	if err != nil {
		log.Fatal(err)
	}
	
	// Downloads file from a remote pod to local file system
	err = cli.Download(context.TODO(), "podname", "containername", "namespace", "testdata/upload.txt", "testdata")
	if err != nil {
		log.Fatal(err)
	}
}
```

## Documentation

You can find the docs at [go docs](https://pkg.go.dev/github.com/shipengqi/kube).

## Test

```bash
go test -v . -kubeconfig <kubeconfig file>
```

Tests Client.Upload, Client.Download And Client.Exec with specified pod, container, namespace:

```bash
go test -v -kubeconfig <kubeconfig file> -container <container name> -pod <pod name> -namespace <namespace> .
```

## 🔋 JetBrains OS licenses

`kube` had been being developed with **GoLand** under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my thanks here.

<a href="https://www.jetbrains.com/?from=kube" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg" alt="JetBrains Logo (Main) logo." width="250" align="middle"></a>