package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_defaultTestPodName       = "nginx-for-exec-test"
	_defaultTestContainerName = "nginx"
	_defaultTestNamespace     = "default"
)

func TestClientExec(t *testing.T) {
	_, err := mockcli.GetPod(context.TODO(), _defaultTestNamespace, _defaultTestPodName)
	if err != nil {
		// pod not found, apply the default test file
		err = mockcli.Apply("testdata/pod-apply.yaml")
		require.NoError(t, err)
		got, err := mockcli.GetPod(context.TODO(), _defaultTestNamespace, _defaultTestPodName)
		t.Log(err)
		t.Log(got)
	}

	stdout, _, err := mockcli.Exec(_defaultTestPodName, _defaultTestContainerName, _defaultTestNamespace, "ls /")
	require.NoError(t, err)
	t.Log(stdout)
}
