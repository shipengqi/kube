package kube

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_defaultTestNamespace = "kube-system"
	_defaultTestPod       = "coredns"
)

func TestClientExec(t *testing.T) {

	containerName := "nginx"
	podNamespace := "default"
	podName := "nginx"

	err := mockcli.Apply("testdata/pod-apply.yaml")
	require.NoError(t, err)

	list, err := mockcli.GetPods(context.TODO(), podNamespace)
	require.NoError(t, err)
	require.NotEqual(t, len(list.Items), 0)

	for _, v := range list.Items {
		if strings.Contains(strings.ToLower(v.Name), podName) {
			podName = v.Name
			containerName = v.Spec.Containers[0].Name
			break
		}
	}
	for i := 0; i < 10; i++ {
		p, err := mockcli.GetPod(context.TODO(), podNamespace, podName)
		require.NoError(t, err)
		t.Log("check pod status:", podName, podNamespace, p.Status, "=====", p.Status.Message)
	}

	// Todo, create a pod, and test Exec(), currently just skip validating the error
	t.Log("exec:", podName, containerName)
	stdout, _, err := mockcli.Exec(podName, containerName, podNamespace, "/bin/sh", "-c", "/bin/ls /")
	t.Log(err.Error())
	require.Error(t, err)
	t.Log(stdout)
}
