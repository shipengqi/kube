package kube

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientExec(t *testing.T) {
	if len(podname) == 0 || len(containername) == 0 || len(namespace) == 0 {
		return
	}
	stdout, _, err := mockcli.Exec(podname, containername, namespace, "ls /")
	require.NoError(t, err)
	t.Log(stdout)
}
