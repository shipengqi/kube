package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientNode(t *testing.T) {
	list, err := mockcli.GetNodes(context.TODO())
	require.NoError(t, err)
	if len(list.Items) == 0 {
		t.Fatal(err)
	}
	node, err := mockcli.GetNode(context.TODO(), list.Items[0].Name)
	require.NoError(t, err)
	assert.Equal(t, list.Items[0].Namespace, node.Namespace)
	assert.Equal(t, list.Items[0].Name, node.Name)
}

func TestClientNamespace(t *testing.T) {
	mocknsname := "mocknamespace"
	t.Run("get:namespace", func(t *testing.T) {
		ns, err := mockcli.GetNamespace(context.TODO(), "default")
		require.NoError(t, err)
		assert.Equal(t, "default", ns.Name)
	})

	t.Run("get:namespace:list", func(t *testing.T) {
		list, err := mockcli.GetNamespaces(context.TODO())
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("create:namespace", func(t *testing.T) {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: mocknsname,
			},
		}
		_, err := mockcli.CreateNamespace(context.TODO(), ns)
		require.NoError(t, err)

		got, err := mockcli.GetNamespace(context.TODO(), mocknsname)
		require.NoError(t, err)
		assert.Equal(t, mocknsname, got.Name)
	})

	t.Run("delete:namespace", func(t *testing.T) {
		err := mockcli.DeleteNamespace(context.TODO(), mocknsname)
		require.NoError(t, err)
	})
}
