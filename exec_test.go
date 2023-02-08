package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_defaultTestNamespace = "kube-system"
)

func TestClientExec(t *testing.T) {
	got, err := mockcli.GetPods(context.TODO(), _defaultTestNamespace)
	require.NoError(t, err)
	require.NotEqual(t, len(got.Items), 0)

	firstPod := got.Items[0]
	firstPodName := firstPod.Name
	firstContainerName := firstPod.Spec.Containers[0].Name

	stdout, _, err := mockcli.Exec(firstPodName, firstContainerName, _defaultTestNamespace, "ls /")
	require.NoError(t, err)
	t.Log(stdout)
}
