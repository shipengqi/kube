package kube

import (
	"context"

	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetIngress returns an Ingress with the given name.
func (c *Client) GetIngress(ctx context.Context, namespace, name string) (*networkv1.Ingress, error) {
	return c.client.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetIngresses returns a IngressList.
func (c *Client) GetIngresses(ctx context.Context, namespace string, label ...string) (*networkv1.IngressList, error) {
	return c.client.NetworkingV1().Ingresses(namespace).List(ctx, listOptions(label))
}

// CreateIngress creates a new Ingress.
func (c *Client) CreateIngress(ctx context.Context, ing *networkv1.Ingress) (*networkv1.Ingress, error) {
	if len(ing.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.NetworkingV1().Ingresses(ing.Namespace).Create(ctx, ing, metav1.CreateOptions{})
}

// DeleteIngress deletes a Ingress.
func (c *Client) DeleteIngress(ctx context.Context, namespace, name string) error {
	return c.client.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetIngressClass returns an IngressClass with the given name.
func (c *Client) GetIngressClass(ctx context.Context, name string) (*networkv1.IngressClass, error) {
	return c.client.NetworkingV1().IngressClasses().Get(ctx, name, metav1.GetOptions{})
}

// GetIngressClasses returns a IngressClassList.
func (c *Client) GetIngressClasses(ctx context.Context, label ...string) (*networkv1.IngressClassList, error) {
	return c.client.NetworkingV1().IngressClasses().List(ctx, listOptions(label))
}

// CreateIngressClass creates a new IngressClass.
func (c *Client) CreateIngressClass(ctx context.Context, ingc *networkv1.IngressClass) (*networkv1.IngressClass, error) {
	return c.client.NetworkingV1().IngressClasses().Create(ctx, ingc, metav1.CreateOptions{})
}

// DeleteIngressClass deletes a IngressClass.
func (c *Client) DeleteIngressClass(ctx context.Context, name string) error {
	return c.client.NetworkingV1().IngressClasses().Delete(ctx, name, metav1.DeleteOptions{})
}
