package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetServices returns a ServiceList.
func (c *Client) GetServices(ctx context.Context, namespace string, label ...string) (*corev1.ServiceList, error) {
	return c.client.CoreV1().Services(namespace).List(ctx, listOptions(label))
}

// GetService returns a Service with the given name.
func (c *Client) GetService(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	return c.client.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateService creates a new Service.
func (c *Client) CreateService(ctx context.Context, svc *corev1.Service) (*corev1.Service, error) {
	if len(svc.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().Services(svc.Namespace).Create(ctx, svc, metav1.CreateOptions{})
}

// DeleteService deletes a Service.
func (c *Client) DeleteService(ctx context.Context, namespace, name string) error {
	return c.client.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
