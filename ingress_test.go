package kube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientIngress(t *testing.T) {
	mockingressname := "mockingress"
	mockingressnamespace := "default"

	// clean
	if _, err := mockcli.GetIngress(context.TODO(), mockingressnamespace, mockingressname); err == nil {
		err := mockcli.DeleteIngress(context.TODO(), mockingressnamespace, mockingressname)
		require.NoError(t, err)
	}

	t.Run("create:ingress", func(t *testing.T) {
		ty := networkv1.PathType("Prefix")
		ing := &networkv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mockingressname,
				Namespace: mockingressnamespace,
			},
			Spec: networkv1.IngressSpec{
				Rules: []networkv1.IngressRule{{
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{Paths: []networkv1.HTTPIngressPath{{
							Path:     "/mock",
							PathType: &ty,
							Backend:  networkv1.IngressBackend{
								Service: &networkv1.IngressServiceBackend{
									Name: "mocksvc",
									Port: networkv1.ServiceBackendPort{
										Number: 8001,
									},
								},
							},
						}}},
					},
				}},
			},
		}
		_, err := mockcli.CreateIngress(context.TODO(), ing)
		require.NoError(t, err)
	})

	t.Run("get:ingress", func(t *testing.T) {
		_, err := mockcli.GetIngress(context.TODO(), mockingressnamespace, mockingressname)
		require.NoError(t, err)
	})

	t.Run("get:ingress:list", func(t *testing.T) {
		list, err := mockcli.GetIngresses(context.TODO(), mockingressnamespace)
		require.NoError(t, err)
		if len(list.Items) == 0 {
			t.Fatal(err)
		}
	})

	t.Run("delete:ingress", func(t *testing.T) {
		err := mockcli.DeleteIngress(context.TODO(), mockingressnamespace, mockingressname)
		require.NoError(t, err)
	})
}
