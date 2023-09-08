package kube

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
)

type ExecRequestMethod string

const (
	ExecPostRequest ExecRequestMethod = "POST"
	ExecGetRequest  ExecRequestMethod = "GET"
)

// Exec is like kubectl exec.
func (c *Client) Exec(pod, container, namespace string, command ...string) (string, string, error) {
	return c.exec(context.TODO(), pod, container, namespace, command...)
}

// ExecWithContext executes exec with context.
func (c *Client) ExecWithContext(ctx context.Context, pod, container, namespace string, command ...string) (string, string, error) {
	return c.exec(ctx, pod, container, namespace, command...)
}

func (c *Client) exec(ctx context.Context, pod, container, namespace string, command ...string) (string, string, error) {
	err := c.validateExecResource(ctx, pod, container, namespace)
	if err != nil {
		return "", "", err
	}

	req := c.RemoteExecRequest(ExecPostRequest, pod, namespace, &corev1.PodExecOptions{
		TypeMeta:  metav1.TypeMeta{},
		Stdout:    true,
		Stderr:    true,
		Container: container,
		Command:   command,
	})
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
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func (c *Client) validateExecResource(ctx context.Context, pod, container, namespace string) error {
	p, err := c.GetPod(ctx, namespace, pod)
	if err != nil {
		return err
	}
	if container != "" {
		found := false
		for _, co := range p.Spec.Containers {
			if co.Name == container {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s not found", container)
		}
	}
	return nil
}
