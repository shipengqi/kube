package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetSecrets returns a SecretList.
func (c *Client) GetSecrets(ctx context.Context, namespace string, label ...string) (*corev1.SecretList, error) {
	return c.client.CoreV1().Secrets(namespace).List(ctx, listOptions(label))
}

// GetSecret returns a Secret with the given name.
func (c *Client) GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	return c.client.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateSecret creates a new Secret.
func (c *Client) CreateSecret(ctx context.Context, secret *corev1.Secret) (*corev1.Secret, error) {
	if len(secret.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
}

// UpdateSecret updates the given Secret.
func (c *Client) UpdateSecret(ctx context.Context, secret *corev1.Secret) (*corev1.Secret, error) {
	if len(secret.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
}

// PatchSecret patch the given Secret.
func (c *Client) PatchSecret(ctx context.Context, namespace, name string, data []byte) (*corev1.Secret, error) {
	return c.client.CoreV1().Secrets(namespace).Patch(
		ctx,
		name,
		types.StrategicMergePatchType,
		data,
		metav1.PatchOptions{},
	)
}

// ApplySecret updates a Secret, and creates a new Secret if not exist.
func (c *Client) ApplySecret(ctx context.Context, namespace, name string, data map[string]string) (*corev1.Secret, error) {
	olds, err := c.GetSecret(ctx, namespace, name)
	if err != nil || olds == nil {
		var news corev1.Secret
		news.SetName(name)
		news.SetNamespace(namespace)
		news.StringData = data
		return c.CreateSecret(ctx, &news)
	}
	olds.StringData = data
	return c.UpdateSecret(ctx, olds)
}

// ApplySecretBytes is like ApplySecret, but with data map[string][]byte.
func (c *Client) ApplySecretBytes(ctx context.Context, namespace, name string, data map[string][]byte) (*corev1.Secret, error) {
	olds, err := c.GetSecret(ctx, namespace, name)
	if err != nil || olds == nil {
		var news corev1.Secret
		news.SetName(name)
		news.SetNamespace(namespace)
		news.Data = data
		return c.CreateSecret(ctx, &news)
	}
	olds.Data = data
	return c.UpdateSecret(ctx, olds)
}

// DeleteSecret deletes a Secret.
func (c *Client) DeleteSecret(ctx context.Context, namespace, name string) error {
	return c.client.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
