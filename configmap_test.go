package kube

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientConfigmap(t *testing.T) {
	mockcmname := "mockconfigmap"
	mockcmnamespace := "default"
	mockcmlabelkey := "mock-configmap"
	mockcmlabelval := "true"
	mockcmdata := map[string]string{
		"mockdata": "mockdata",
	}

	// clean
	if _, err := mockcli.GetConfigMap(context.TODO(), mockcmnamespace, mockcmname); err == nil {
		err := mockcli.DeleteConfigMap(context.TODO(), mockcmnamespace, mockcmname)
		require.NoError(t, err)
	}

	t.Run("create:configmap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mockcmname,
				Namespace: mockcmnamespace,
				Labels: map[string]string{
					mockcmlabelkey: mockcmlabelval,
				},
			},
			Data: mockcmdata,
		}
		_, err := mockcli.CreateConfigMap(context.TODO(), cm)
		require.NoError(t, err)
	})

	t.Run("create:configmap:error", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: mockcmname,
			},
			Data: mockcmdata,
		}
		_, err := mockcli.CreateConfigMap(context.TODO(), cm)
		require.EqualError(t, err, ErrorMissingNamespace.Error())
	})

	t.Run("get:configmap", func(t *testing.T) {
		cm, err := mockcli.GetConfigMap(context.TODO(), mockcmnamespace, mockcmname)
		require.NoError(t, err)
		assert.Equal(t, mockcmdata, cm.Data)
	})

	t.Run("get:configmap:list", func(t *testing.T) {
		list, err := mockcli.GetConfigMaps(context.TODO(), mockcmnamespace)
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("get:configmap:list:label", func(t *testing.T) {
		list, err := mockcli.GetConfigMaps(context.TODO(), mockcmnamespace,
			fmt.Sprintf("%s=%s", mockcmlabelkey, mockcmlabelval))
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:configmap", func(t *testing.T) {
		err := mockcli.DeleteConfigMap(context.TODO(), mockcmnamespace, mockcmname)
		require.NoError(t, err)
	})
}
