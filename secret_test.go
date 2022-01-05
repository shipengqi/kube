package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientSecret(t *testing.T) {
	mocksecretname := "mocksecret"
	mocksecretnamespace := "default"
	mocksecretdata := map[string][]byte{
		"mockdata": []byte("mockdata"),
	}

	// clean
	if _, err := mockcli.GetSecret(context.TODO(), mocksecretnamespace, mocksecretname); err == nil {
		err := mockcli.DeleteSecret(context.TODO(), mocksecretnamespace, mocksecretname)
		require.NoError(t, err)
	}

	t.Run("create:secret", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mocksecretname,
				Namespace: mocksecretnamespace,
			},
			Data: mocksecretdata,
		}
		_, err := mockcli.CreateSecret(context.TODO(), secret)
		require.NoError(t, err)
	})

	t.Run("create:secret:error", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: mocksecretname,
			},
			Data: mocksecretdata,
		}
		_, err := mockcli.CreateSecret(context.TODO(), secret)
		require.EqualError(t, err, ErrorMissingNamespace.Error())
	})

	t.Run("get:secret", func(t *testing.T) {
		secret, err := mockcli.GetSecret(context.TODO(), mocksecretnamespace, mocksecretname)
		require.NoError(t, err)
		assert.Equal(t, mocksecretdata, secret.Data)
	})

	t.Run("get:secret:list", func(t *testing.T) {
		list, err := mockcli.GetSecrets(context.TODO(), mocksecretnamespace)
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:secret", func(t *testing.T) {
		err := mockcli.DeleteSecret(context.TODO(), mocksecretnamespace, mocksecretname)
		require.NoError(t, err)
	})
}
