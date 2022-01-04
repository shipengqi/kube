package kube

import (
	"bytes"
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// Exec is like kubectl exec
func (c *Client) Exec(pod, container, namespace, command string) (string, string, error) {
	p, err := c.GetPod(context.TODO(), namespace, pod)
	if err != nil {
		return "", "", err
	}
	req := c.client.CoreV1().RESTClient().Post().
		Resource("pods").Name(p.Name).
		Namespace(p.Namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		TypeMeta:  metav1.TypeMeta{},
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
		Container: container,
		Command:   []string{"sh", "-c", command},
	}, scheme.ParameterCodec)

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	restcfg, err := c.RestConfig()
	if err != nil {
		return "", "", err
	}
	exec, err := remotecommand.NewSPDYExecutor(restcfg, "POST", req.URL())
	if err != nil {
		return "", "", err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}
