package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientRole(t *testing.T)  {
	mockrolename := "mockrole"
	mockrolenamespace := "default"

	// clean
	if _, err := mockcli.GetRole(context.TODO(), mockrolenamespace, mockrolename); err == nil {
		err := mockcli.DeleteRole(context.TODO(), mockrolenamespace, mockrolename)
		require.NoError(t, err)
	}

	t.Run("create:role", func(t *testing.T) {
		r := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name: mockrolename,
				Namespace: mockrolenamespace,
			},
		}
		_, err := mockcli.CreateRole(context.TODO(), r)
		require.NoError(t, err)
	})

	t.Run("get:role", func(t *testing.T) {
		_, err := mockcli.GetRole(context.TODO(), mockrolenamespace, mockrolename)
		require.NoError(t, err)
	})

	t.Run("get:role:list", func(t *testing.T) {
		list, err := mockcli.GetRoles(context.TODO(), mockrolenamespace)
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:role", func(t *testing.T) {
		err := mockcli.DeleteRole(context.TODO(), mockrolenamespace, mockrolename)
		require.NoError(t, err)
	})
}

func TestClientServiceAccount(t *testing.T)  {
	mocksaname := "mocksa"
	mocksanamespace := "default"

	// clean
	if _, err := mockcli.GetServiceAccount(context.TODO(), mocksanamespace, mocksaname); err == nil {
		err := mockcli.DeleteServiceAccount(context.TODO(), mocksanamespace, mocksaname)
		require.NoError(t, err)
	}

	t.Run("create:serviceaccount", func(t *testing.T) {
		sa := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: mocksaname,
				Namespace: mocksanamespace,
			},
		}
		_, err := mockcli.CreateServiceAccount(context.TODO(), sa)
		require.NoError(t, err)
	})

	t.Run("get:serviceaccount", func(t *testing.T) {
		_, err := mockcli.GetServiceAccount(context.TODO(), mocksanamespace, mocksaname)
		require.NoError(t, err)
	})

	t.Run("get:serviceaccount:list", func(t *testing.T) {
		list, err := mockcli.GetServiceAccounts(context.TODO(), mocksanamespace)
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:serviceaccount", func(t *testing.T) {
		err := mockcli.DeleteServiceAccount(context.TODO(), mocksanamespace, mocksaname)
		require.NoError(t, err)
	})
}

func TestClientClusterRole(t *testing.T)  {
	mockclusterrolename := "mockclusterrole"

	// clean
	if _, err := mockcli.GetClusterRole(context.TODO(), mockclusterrolename); err == nil {
		err := mockcli.DeleteClusterRole(context.TODO(), mockclusterrolename)
		require.NoError(t, err)
	}

	t.Run("create:clusterrole", func(t *testing.T) {
		r := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: mockclusterrolename,
			},
		}
		_, err := mockcli.CreateClusterRole(context.TODO(), r)
		require.NoError(t, err)
	})

	t.Run("get:clusterrole", func(t *testing.T) {
		_, err := mockcli.GetClusterRole(context.TODO(), mockclusterrolename)
		require.NoError(t, err)
	})

	t.Run("get:clusterrole:list", func(t *testing.T) {
		list, err := mockcli.GetClusterRoles(context.TODO())
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:clusterrole", func(t *testing.T) {
		err := mockcli.DeleteClusterRole(context.TODO(), mockclusterrolename)
		require.NoError(t, err)
	})
}
