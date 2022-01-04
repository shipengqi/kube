package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetConfigMaps returns a ConfigMapList.
func (c *Client) GetConfigMaps(ctx context.Context, namespace string, label ...string) (*corev1.ConfigMapList, error) {
	return c.client.CoreV1().ConfigMaps(namespace).List(ctx, listOptions(label))
}

// GetConfigMap returns a ConfigMap with the given name.
func (c *Client) GetConfigMap(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	return c.client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateConfigMap creates a new ConfigMap.
func (c *Client) CreateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if len(cm.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().ConfigMaps(cm.Namespace).Create(ctx, cm, metav1.CreateOptions{})
}

// UpdateConfigMap updates a ConfigMap.
func (c *Client) UpdateConfigMap(ctx context.Context, cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	if len(cm.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().ConfigMaps(cm.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
}

// PatchConfigMap patch a ConfigMap.
func (c *Client) PatchConfigMap(ctx context.Context, namespace, name string, data []byte) (*corev1.ConfigMap, error) {
	return c.client.CoreV1().ConfigMaps(namespace).Patch(
		ctx,
		name,
		types.StrategicMergePatchType,
		data,
		metav1.PatchOptions{},
	)
}

// ApplyConfigMap updates a ConfigMap, and creates a new ConfigMap if not exist.
func (c *Client) ApplyConfigMap(ctx context.Context, namespace, name string, data map[string]string) (*corev1.ConfigMap, error) {
	olds, err := c.GetConfigMap(ctx, namespace, name)
	if err != nil || olds == nil {
		var news corev1.ConfigMap
		news.SetName(name)
		news.SetNamespace(namespace)
		news.Data = data
		return c.CreateConfigMap(ctx, &news)
	}
	olds.Data = data
	return c.UpdateConfigMap(ctx, olds)
}

// DeleteConfigMap deletes a ConfigMap.
func (c *Client) DeleteConfigMap(ctx context.Context, namespace, name string) error {
	return c.client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
