package kube

import (
	"testing"
)

func TestClientExec(t *testing.T) {
	if len(podname) == 0 || len(containername) == 0 || len(namespace) == 0 {
		return
	}

	// _, err := mockcli.GetPod(context.TODO(), namespace, podname)
	// if err != nil {
	// 	// pod not found, apply the default test file
	// 	err = mockcli.Apply("testdata/pod-apply.yaml")
	// 	require.NoError(t, err)
	// }
	//
	// stdout, _, err := mockcli.Exec(podname, containername, namespace, "ls /")
	// require.NoError(t, err)
	// t.Log(stdout)
}
