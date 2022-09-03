package kube

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientConfigmap(t *testing.T) {
	mockcmname := "mockconfigmap"
	mockcmnamenoexist := "mockcmnamenoexist"
	mockcmnamespace := "default"
	mockcmlabelkey := "mock-configmap"
	mockcmlabelval := "true"
	mockcmdata := map[string]string{
		"mockdata": "mockdata",
	}
	mockcmdataupdate := map[string]string{
		"mockdata":       "mockdata",
		"mockdataupdate": "mockdataupdate",
	}
	// clean
	if _, err := mockcli.GetConfigMap(context.TODO(), mockcmnamespace, mockcmname); err == nil {
		err := mockcli.DeleteConfigMap(context.TODO(), mockcmnamespace, mockcmname)
		require.NoError(t, err)
	}
	if _, err := mockcli.GetConfigMap(context.TODO(), mockcmnamespace, mockcmnamenoexist); err == nil {
		err := mockcli.DeleteConfigMap(context.TODO(), mockcmnamespace, mockcmnamenoexist)
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

	t.Run("update:configmap", func(t *testing.T) {
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mockcmname,
				Namespace: mockcmnamespace,
			},
			Data: mockcmdataupdate,
		}
		updated, err := mockcli.UpdateConfigMap(context.TODO(), cm)
		require.NoError(t, err)
		assert.Equal(t, mockcmdataupdate, updated.Data)
	})

	t.Run("apply:configmap:exist", func(t *testing.T) {
		applydata := map[string]string{
			"mockdata":      "mockdata",
			"mockdataapply": "mockdataapply",
		}
		applyed, err := mockcli.ApplyConfigMap(context.TODO(), mockcmnamespace, mockcmname, applydata)
		require.NoError(t, err)
		assert.Equal(t, applydata, applyed.Data)
	})

	t.Run("apply:configmap:no:exist", func(t *testing.T) {
		_, err := mockcli.ApplyConfigMap(context.TODO(), mockcmnamespace, mockcmnamenoexist, mockcmdata)
		require.NoError(t, err)
	})

	t.Run("delete:configmap", func(t *testing.T) {
		err := mockcli.DeleteConfigMap(context.TODO(), mockcmnamespace, mockcmname)
		require.NoError(t, err)
		err = mockcli.DeleteConfigMap(context.TODO(), mockcmnamespace, mockcmnamenoexist)
		require.NoError(t, err)
	})
}
