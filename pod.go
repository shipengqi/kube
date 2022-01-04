package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods returns a PodList.
func (c *Client) GetPods(ctx context.Context, namespace string, label ...string) (*corev1.PodList, error) {
	return c.client.CoreV1().Pods(namespace).List(ctx, listOptions(label))
}

// GetPod returns a Pod with the given name.
func (c *Client) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	return c.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreatePod creates a new Pod.
func (c *Client) CreatePod(ctx context.Context, pod *corev1.Pod) (*corev1.Pod, error) {
	if len(pod.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
}

// DeletePod deletes a Pod.
func (c *Client) DeletePod(ctx context.Context, namespace, name string) error {
	return c.client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
