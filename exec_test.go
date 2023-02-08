package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_defaultTestPodName       = "nginx-for-exec-test"
	_defaultTestContainerName = "nginx"
	_defaultTestNamespace     = "kube-system"
)

func TestClientExec(t *testing.T) {
	got, err := mockcli.GetPods(context.TODO(), _defaultTestNamespace)
	require.NoError(t, err)
	require.NotEqual(t, len(got.Items), 0)
	t.Log(got.Items)

	// stdout, _, err := mockcli.Exec(_defaultTestPodName, _defaultTestContainerName, _defaultTestNamespace, "ls /")
	// require.NoError(t, err)
	// t.Log(stdout)
}
