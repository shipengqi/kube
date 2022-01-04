package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestClientService(t *testing.T) {
	mockservicename := "mockservice"
	mockservicenamespace := "default"
	mockserviceportname := "mockportname"
	var mockserviceport int32 = 8001
	// clean
	if _, err := mockcli.GetService(context.TODO(), mockservicenamespace, mockservicename); err == nil {
		err := mockcli.DeleteService(context.TODO(), mockservicenamespace, mockservicename)
		require.NoError(t, err)
	}

	t.Run("create:service", func(t *testing.T) {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: mockservicename,
				Namespace: mockservicenamespace,
			},
			Spec:       corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:        mockserviceportname,
					Port:        mockserviceport,
					TargetPort:  intstr.IntOrString{
						IntVal: mockserviceport,
					},
				}},
			},
		}
		_, err := mockcli.CreateService(context.TODO(), svc)
		require.NoError(t, err)
	})

	t.Run("get:service", func(t *testing.T) {
		_, err := mockcli.GetService(context.TODO(), mockservicenamespace, mockservicename)
		require.NoError(t, err)
	})

	t.Run("get:service:list", func(t *testing.T) {
		list, err := mockcli.GetServices(context.TODO(), mockservicenamespace)
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:service", func(t *testing.T) {
		err := mockcli.DeleteService(context.TODO(), mockservicenamespace, mockservicename)
		require.NoError(t, err)
	})
}
