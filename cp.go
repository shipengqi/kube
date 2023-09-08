package kube

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
)

// Download downloads file from a remote pod to local file system.
// Requires that the 'tar' binary is present in your container
// image.  If 'tar' is not present, 'Download' will fail.
func (c *Client) Download(ctx context.Context, pod, container, namespace, src, dst string) error {
	if err := c.validateExecResource(ctx, pod, container, namespace); err != nil {
		return err
	}

	command := []string{"tar", "-cf", "-", src}
	req := c.RemoteExecRequest(ExecGetRequest, pod, namespace, &corev1.PodExecOptions{
		TypeMeta:  metav1.TypeMeta{},
		Stdout:    true,
		Stderr:    true,
		Container: container,
		Command:   command,
	})

	exec, err := c.RemoteExecutor(req)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	var ebuf bytes.Buffer
	if err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &buf,
		Stderr: &ebuf,
	}); err != nil {
		return fmt.Errorf("exec.StreamWithContext, %s, %s", err.Error(), ebuf.String())
	}

	prefix := getPrefix(src)
	prefix = path.Clean(prefix)
	prefix = stripPathShortcuts(prefix)
	dstPath := path.Join(dst, path.Base(prefix))
	err = untar(&buf, dstPath, prefix)
	if err != nil {
		return err
	}
	return nil
}

// Upload uploads local file to a remote pod.
// Requires that the 'tar' binary is present in your container
// image.  If 'tar' is not present, 'Upload' will fail.
func (c *Client) Upload(ctx context.Context, pod, container, namespace, src, dst string) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("%s doesn't exist in local filesystem", src)
	}
	if err := c.validateExecResource(ctx, pod, container, namespace); err != nil {
		return err
	}

	tarName := getTmpTarName(src)
	if err := makeTar(src, tarName); err != nil {
		return err
	}
	defer func() { _ = os.Remove(tarName) }()

	if dst != "/" && strings.HasSuffix(string(dst[len(dst)-1]), "/") {
		dst = dst[:len(dst)-1]
	}
	if err := c.checkRemoteDstIsDir(ctx, pod, container, namespace, dst); err == nil {
		dst = dst + "/" + path.Base(src)
	}
	dstDir := path.Dir(dst)

	command := []string{"tar", "-xmf", "-"}
	if len(dstDir) > 0 {
		command = append(command, "-C", dstDir)
	}
	req := c.RemoteExecRequest(ExecGetRequest, pod, namespace, &corev1.PodExecOptions{
		TypeMeta:  metav1.TypeMeta{},
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		Container: container,
		Command:   command,
	})

	exec, err := c.RemoteExecutor(req)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(tarName)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	var ebuf bytes.Buffer
	if err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  bytes.NewReader(data),
		Stdout: &buf,
		Stderr: &ebuf,
	}); err != nil {
		return fmt.Errorf("exec.StreamWithContext, %s, %s", err.Error(), ebuf.String())
	}

	return nil
}

// checkDestinationIsDir receives a destination string and
// determines if the provided destination path exists on the
// pod. If the destination path does not exist or is not a
// directory, an error is returned with the exit code received.
func (c *Client) checkRemoteDstIsDir(ctx context.Context, pod, container, namespace, dst string) error {
	command := []string{"test", "-d", dst}
	req := c.RemoteExecRequest(ExecGetRequest, pod, namespace, &corev1.PodExecOptions{
		TypeMeta:  metav1.TypeMeta{},
		Stdout:    true,
		Container: container,
		Command:   command,
	})

	exec, err := c.RemoteExecutor(req)
	if err != nil {
		return err
	}

	return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: bytes.NewBuffer([]byte{}),
	})
}

func makeTar(src, dst string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = fw.Close() }()

	return tarf(fw, src)
}

func tarf(writer io.Writer, src string) error {
	tw := tar.NewWriter(writer)
	defer func() { _ = tw.Close() }()
	return filepath.Walk(src, func(filename string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(filename, string(filepath.Separator))
		// write file info
		if err = tw.WriteHeader(header); err != nil {
			return err
		}
		// whether info describes a regular file.
		if !info.Mode().IsRegular() {
			return nil
		}
		fr, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer func() { _ = fr.Close() }()
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
		return nil
	})
}

func untar(reader io.Reader, dst, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
				return err
			}
			break
		}
		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}
		dstPath := filepath.Join(dst, header.Name[len(prefix):])

		baseName := filepath.Dir(dstPath)
		if err = os.MkdirAll(baseName, os.ModePerm); err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			if err = createFile(tarReader, dstPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}
	return nil
}

func createFile(src io.Reader, dst string, mode os.FileMode) error {
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if _, err = io.Copy(file, src); err != nil {
		return err
	}
	return nil
}

func getPrefix(fpath string) string {
	return strings.TrimLeft(fpath, "/")
}

func getTmpTarName(fpath string) string {
	names := strings.Split(path.Base(fpath), ".")
	return fmt.Sprintf("%s/%s.tar", os.TempDir(), names[0])
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
func stripPathShortcuts(p string) string {
	newPath := p
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}
