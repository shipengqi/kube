package kube

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var containerName string

	podName, containerName, podNamespace := getResourceNames(t)

	t.Run("should upload error", func(t *testing.T) {
		err := mockcli.Upload(context.TODO(), podName, containerName, podNamespace, "testdata/noexists.txt", "/testdata")
		require.ErrorContains(t, err, "doesn't exist in local filesystem")
	})

	t.Run("should upload file to sub dir error", func(t *testing.T) {
		err := mockcli.Upload(context.TODO(), podName, containerName, podNamespace, "testdata/upload.txt", "/testdata/sub/upload.txt")
		require.ErrorContains(t, err, "/testdata/sub: Cannot open: No such file or directory")
	})

	t.Run("upload file to dir", func(t *testing.T) {
		err := mockcli.Upload(context.TODO(), podName, containerName, podNamespace, "testdata/upload.txt", "/testdata")
		require.NoError(t, err)

		out, _, err := mockcli.Exec(podName, containerName, podNamespace, "ls", "/testdata")
		require.NoError(t, err)
		assert.Contains(t, out, "upload.txt")

	})

	t.Run("should download err", func(t *testing.T) {
		err := mockcli.Download(context.TODO(), podName, containerName, podNamespace, "testdata/noexists.txt", "testdata")
		t.Log(err.Error())
	})

	t.Run("download file", func(t *testing.T) {
		err := mockcli.Download(context.TODO(), podName, containerName, podNamespace, "testdata/upload.txt", "testdata")
		require.NoError(t, err)
		files, err := os.ReadDir("testdata")
		require.NoError(t, err)
		for _, v := range files {
			t.Log(v.Name(), v.IsDir())
		}
		_, _, err = mockcli.Exec(podName, containerName, podNamespace, "rm", "-rf", "/testdata")
		require.NoError(t, err)
	})

}
