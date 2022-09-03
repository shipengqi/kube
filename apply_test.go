package kube

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientApply(t *testing.T) {
	k8s, err := mockcli.Dial()
	require.NoError(t, err)

	tests := []struct {
		name  string
		files []string
		sas   []string
		cms   []testcm
	}{
		{
			"create:file:one",
			[]string{"testdata/content-create.yaml"},
			[]string{"default:kube-sa1", "default:kube-sa2"},
			[]testcm{
				{"default:kube-cm1", map[string]string{"TESTDATA": "kube-cm1"}},
			},
		},
		{
			"apply:file:one",
			[]string{"testdata/content-apply.yaml"},
			[]string{"default:kube-sa1", "default:kube-sa2"},
			[]testcm{
				{"default:kube-cm1", map[string]string{"TESTDATA": "kube-cm1-update"}},
				{"default:kube-cm2", map[string]string{"TESTDATA": "kube-cm2"}},
			},
		},
		{
			"apply:file:multi",
			[]string{"testdata/content-apply.yaml", "testdata/content-apply2.yaml"},
			[]string{"default:kube-sa1", "default:kube-sa2"},
			[]testcm{
				{"default:kube-cm1", map[string]string{"TESTDATA": "kube-cm1-multi"}},
				{"default:kube-cm2", map[string]string{"TESTDATA": "kube-cm2", "TESTDATA2": "kube-cm2-multi"}},
			},
		},
	}

	// clean
	err = mockcli.Delete("testdata/content-create.yaml", "testdata/content-apply.yaml", "testdata/content-apply2.yaml")
	require.NoError(t, err)

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			err = mockcli.Apply(v.files...)
			require.NoError(t, err)
			time.Sleep(time.Second * 1)
			for _, sa := range v.sas {
				tokens := strings.Split(sa, ":")
				_, err = k8s.CoreV1().ServiceAccounts(tokens[0]).Get(context.TODO(), tokens[1], metav1.GetOptions{})
				require.NoError(t, err)
				time.Sleep(time.Millisecond * 200)
			}
			for _, cm := range v.cms {
				tokens := strings.Split(cm.name, ":")
				d, err := k8s.CoreV1().ConfigMaps(tokens[0]).Get(context.TODO(), tokens[1], metav1.GetOptions{})
				require.NoError(t, err)
				require.Equal(t, cm.data, d.Data)
				time.Sleep(time.Millisecond * 200)
			}
		})
	}

	// clean
	err = mockcli.Delete("testdata/content-create.yaml", "testdata/content-apply.yaml", "testdata/content-apply2.yaml")
	require.NoError(t, err)
}
