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
	mocksecretnamenoexist := "mocksecretnoexist"
	mocksecretnamespace := "default"
	mocksecretdata := map[string][]byte{
		"mockdata": []byte("mockdata"),
	}

	// clean
	if _, err := mockcli.GetSecret(context.TODO(), mocksecretnamespace, mocksecretname); err == nil {
		err := mockcli.DeleteSecret(context.TODO(), mocksecretnamespace, mocksecretname)
		require.NoError(t, err)
	}
	if _, err := mockcli.GetSecret(context.TODO(), mocksecretnamespace, mocksecretnamenoexist); err == nil {
		err := mockcli.DeleteSecret(context.TODO(), mocksecretnamespace, mocksecretnamenoexist)
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

	t.Run("update:secret", func(t *testing.T) {
		updatedata := map[string][]byte{
			"mockdata":       []byte("mockdata"),
			"mockdataupdate": []byte("mockdata"),
		}
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mocksecretname,
				Namespace: mocksecretnamespace,
			},
			Data: updatedata,
		}
		updated, err := mockcli.UpdateSecret(context.TODO(), secret)
		require.NoError(t, err)
		assert.Equal(t, updatedata, updated.Data)
	})

	t.Run("apply:secret:exist", func(t *testing.T) {
		applydata := map[string]string{
			"mockdata":      "mockdata",
			"mockdataapply": "mockdataapply",
		}
		applyed, err := mockcli.ApplySecret(context.TODO(), mocksecretnamespace, mocksecretname, applydata)
		require.NoError(t, err)
		assert.Equal(t, map[string][]byte{
			"mockdata":       []byte("mockdata"),
			"mockdataapply":  []byte("mockdataapply"),
			"mockdataupdate": []byte("mockdata"),
		}, applyed.Data)
	})

	t.Run("apply:secret:no:exist", func(t *testing.T) {
		applydata := map[string]string{
			"mockdata":      "mockdata",
			"mockdataapply": "mockdataapply",
		}
		_, err := mockcli.ApplySecret(context.TODO(), mocksecretnamespace, mocksecretnamenoexist, applydata)
		require.NoError(t, err)
	})

	t.Run("apply:secret:bytes:exist", func(t *testing.T) {
		applydata := map[string][]byte{
			"mockdata":      []byte("mockdata"),
			"mockdataapply": []byte("mockdataapply"),
		}
		applyed, err := mockcli.ApplySecretBytes(context.TODO(), mocksecretnamespace, mocksecretname, applydata)
		require.NoError(t, err)
		assert.Equal(t, applydata, applyed.Data)
	})

	t.Run("apply:secret:bytes:no:exist", func(t *testing.T) {
		applydata := map[string][]byte{
			"mockdata":      []byte("mockdata"),
			"mockdataapply": []byte("mockdataapply"),
		}
		_, err := mockcli.ApplySecretBytes(context.TODO(), mocksecretnamespace, mocksecretnamenoexist, applydata)
		require.NoError(t, err)
	})

	t.Run("delete:secret", func(t *testing.T) {
		err := mockcli.DeleteSecret(context.TODO(), mocksecretnamespace, mocksecretname)
		require.NoError(t, err)

		err = mockcli.DeleteSecret(context.TODO(), mocksecretnamespace, mocksecretnamenoexist)
		require.NoError(t, err)
	})
}
