package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetServiceAccounts returns a ServiceAccountList.
func (c *Client) GetServiceAccounts(ctx context.Context, namespace string, label ...string) (*corev1.ServiceAccountList, error) {
	return c.client.CoreV1().ServiceAccounts(namespace).List(ctx, listOptions(label))
}

// GetServiceAccount returns a ServiceAccount with the given name.
func (c *Client) GetServiceAccount(ctx context.Context, namespace, name string) (*corev1.ServiceAccount, error) {
	return c.client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateServiceAccount creates a new ServiceAccount.
func (c *Client) CreateServiceAccount(ctx context.Context, sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	if len(sa.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.CoreV1().ServiceAccounts(sa.Namespace).Create(ctx, sa, metav1.CreateOptions{})
}

// DeleteServiceAccount deletes a ServiceAccount.
func (c *Client) DeleteServiceAccount(ctx context.Context, namespace, name string) error {
	return c.client.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetClusterRoles returns a ClusterRoleList.
func (c *Client) GetClusterRoles(ctx context.Context, label ...string) (*rbacv1.ClusterRoleList, error) {
	return c.client.RbacV1().ClusterRoles().List(ctx, listOptions(label))
}

// GetClusterRole returns a ClusterRole with the given name.
func (c *Client) GetClusterRole(ctx context.Context, name string) (*rbacv1.ClusterRole, error) {
	return c.client.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
}

// CreateClusterRole creates a new ClusterRole.
func (c *Client) CreateClusterRole(ctx context.Context, role *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	return c.client.RbacV1().ClusterRoles().Create(ctx, role, metav1.CreateOptions{})
}

// DeleteClusterRole deletes a ClusterRole.
func (c *Client) DeleteClusterRole(ctx context.Context, name string) error {
	return c.client.RbacV1().ClusterRoles().Delete(ctx, name, metav1.DeleteOptions{})
}

// GetRoles returns a RoleList.
func (c *Client) GetRoles(ctx context.Context, namespace string, label ...string) (*rbacv1.RoleList, error) {
	return c.client.RbacV1().Roles(namespace).List(ctx, listOptions(label))
}

// GetRole returns a Role with the given name.
func (c *Client) GetRole(ctx context.Context, namespace, name string) (*rbacv1.Role, error) {
	return c.client.RbacV1().Roles(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateRole create a new Role.
func (c *Client) CreateRole(ctx context.Context, role *rbacv1.Role) (*rbacv1.Role, error) {
	if len(role.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.RbacV1().Roles(role.Namespace).Create(ctx, role, metav1.CreateOptions{})
}

// DeleteRole deletes a Role.
func (c *Client) DeleteRole(ctx context.Context, namespace, name string) error {
	return c.client.RbacV1().Roles(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
