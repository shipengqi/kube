package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
)

// GetVersion returns server's version.
func (c *Client) GetVersion() (*version.Info, error) {
	return c.client.Discovery().ServerVersion()
}

// GetNodes returns a NodeList.
func (c *Client) GetNodes(ctx context.Context, label ...string) (*corev1.NodeList, error) {
	return c.client.CoreV1().Nodes().List(ctx, listOptions(label))
}

// GetNode returns a Node with the given name.
func (c *Client) GetNode(ctx context.Context, name string) (*corev1.Node, error) {
	return c.client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
}

// DeleteNode deletes a Node.
func (c *Client) DeleteNode(ctx context.Context, name string) error {
	return c.client.CoreV1().Nodes().Delete(ctx, name, metav1.DeleteOptions{})
}

// GetNamespaces returns a NamespaceList.
func (c *Client) GetNamespaces(ctx context.Context, label ...string) (*corev1.NamespaceList, error) {
	return c.client.CoreV1().Namespaces().List(ctx, listOptions(label))
}

// GetNamespace returns a Namespace with the given name.
func (c *Client) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	return c.client.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
}

// CreateNamespace creates a Namespace.
func (c *Client) CreateNamespace(ctx context.Context, namespace *corev1.Namespace) (*corev1.Namespace, error) {
	return c.client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
}

// DeleteNamespace deletes a Node.
func (c *Client) DeleteNamespace(ctx context.Context, name string) error {
	return c.client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}
