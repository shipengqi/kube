package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientPod(t *testing.T) {
	mocknamespace := "default"
	mockpodname := "nginx"
	// clean
	if _, err := mockcli.GetPod(context.TODO(), mocknamespace, mockpodname); err == nil {
		err := mockcli.DeletePod(context.TODO(), mocknamespace, mockpodname)
		require.NoError(t, err)
	}

	t.Run("create:pod", func(t *testing.T) {
		p := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      mockpodname,
				Namespace: mocknamespace,
				Labels: map[string]string{
					"app": "nginx",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "nginx",
					},
				},
			},
		}

		_, err := mockcli.CreatePod(context.TODO(), p)
		require.NoError(t, err)
	})

	t.Run("get:pod", func(t *testing.T) {
		_, err := mockcli.GetPod(context.TODO(), mocknamespace, mockpodname)
		require.NoError(t, err)
	})

	t.Run("get:pod:list", func(t *testing.T) {
		list, err := mockcli.GetPods(context.TODO(), mocknamespace)
		require.NoError(t, err)
		require.NotEqual(t, list.Items, 0)
	})

	t.Run("delete:pod", func(t *testing.T) {
		err := mockcli.DeletePod(context.TODO(), mocknamespace, mockpodname)
		require.NoError(t, err)
	})
}
