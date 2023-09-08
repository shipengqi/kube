package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var containerName string

	podName, containerName, podNamespace := getResourceNames(t)

	t.Run("upload", func(t *testing.T) {
		err := mockcli.Upload(context.TODO(), podName, containerName, podNamespace, "testdata/upload.txt", "/testdata")
		require.NoError(t, err)
	})

	t.Run("download", func(t *testing.T) {
		err := mockcli.Download(context.TODO(), podName, containerName, podNamespace, "testdata/upload.txt", "testdata")
		require.NoError(t, err)
	})

}
