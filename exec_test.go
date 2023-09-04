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

	// exec coredns
	for _, v := range list.Items {
		if strings.Contains(strings.ToLower(v.Name), podName) {
			podName = v.Name
			containerName = v.Spec.Containers[0].Name
			break
		}
	}

	// Todo, create a pod, and test Exec(), currently just skip validating the error
	t.Log("exec:", "nginx", "nginx")
	stdout, _, err := mockcli.Exec(podName, containerName, podNamespace, "/bin/sh", "-c", "/bin/ls /")
	t.Log(err.Error())
	require.Error(t, err)
	t.Log(stdout)
}
