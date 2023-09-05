package kube

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_defaultTestNamespace = "ingress-nginx"
	_defaultTestPod       = "nginx"
)

func TestClientExec(t *testing.T) {

	var containerName string

	podNamespace := _defaultTestNamespace
	podName := _defaultTestPod

	if namespace != "" {
		podNamespace = namespace
	}

	if podname != "" && containername != "" {
		podName = podname
		containerName = containername
	} else {
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
	}

	t.Log("exec:", podName, containerName)
	stdout, _, err := mockcli.Exec(podName, containerName, podNamespace, "echo", "hello")
	require.NoError(t, err)
	t.Log(stdout)
	assert.Equal(t, "hello", stdout)
}
