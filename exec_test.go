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
	list, err := mockcli.GetPods(context.TODO(), _defaultTestNamespace)
	require.NoError(t, err)
	require.NotEqual(t, len(list.Items), 0)

	var (
		podName       string
		containerName string
	)
	// exec -it coredns
	for _, v := range list.Items {
		if strings.Contains(strings.ToLower(v.Name), _defaultTestPod) {
			podName = v.Name
			containerName = v.Spec.Containers[0].Name
		}
	}

	t.Log("exec:", podName, containerName)
	stdout, _, err := mockcli.Exec(podName, containerName, _defaultTestNamespace, "/bin/ls /")
	require.NoError(t, err)
	t.Log(stdout)
}
