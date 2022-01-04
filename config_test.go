package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type currentFuncRes struct {
	err   string
	value string
}

func TestConfigContexts(t *testing.T) {
	tests := []struct {
		file    string
		ctx     string
		user    currentFuncRes
		cluster currentFuncRes
	}{
		{
			"testdata/config", "dev-frontend",
			currentFuncRes{"", "developer"},
			currentFuncRes{"", "development"},
		},
		{
			"testdata/multi-ctx-config", "",
			currentFuncRes{"not found", ""},
			currentFuncRes{"not found", ""},
		},
	}

	for k := range tests {
		c := tests[k]
		t.Run(c.file, func(t *testing.T) {
			flags := genericclioptions.NewConfigFlags(false)
			flags.KubeConfig = &c.file
			cfg := NewConfig(flags)
			got, err := cfg.CurrentContextName()
			assert.NoError(t, err)
			assert.Equal(t, c.ctx, got)

			got, err = cfg.CurrentClusterName()
			if c.cluster.err != "" {
				assert.Contains(t, err.Error(), c.cluster.err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, c.cluster.value, got)

			got, err = cfg.CurrentUserName()
			if c.user.err != "" {
				assert.Contains(t, err.Error(), c.user.err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, c.user.value, got)
		})
	}
}

func TestConfigUseContext(t *testing.T) {
	c := "testdata/multi-ctx-config"
	ctx := "dev-frontend"
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &c
	cfg := NewConfig(flags)

	// reset ctx
	err := cfg.useContext("")
	assert.NoError(t, err)

	err = cfg.UseContext(ctx)
	assert.NoError(t, err)

	got, err := cfg.CurrentContextName()
	assert.Equal(t, ctx, got)

	err = cfg.UseContext("notfound")
	assert.Contains(t, err.Error(), "not found")

	// clean ctx
	err = cfg.useContext("")
	assert.NoError(t, err)
}

func TestConfigSetContext(t *testing.T) {
	c := "testdata/config"
	ctxname := "testset"
	ctx := &clientcmdapi.Context{
		Cluster:   "testsetcluster",
		AuthInfo:  "testsetuser",
		Namespace: "testsetnamespace",
	}
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &c
	cfg := NewConfig(flags)

	err := cfg.SetContext(ctxname, ctx)
	assert.NoError(t, err)

	got, err := cfg.GetContext(ctxname)
	assert.NoError(t, err)

	assert.Equal(t, ctx.AuthInfo, got.AuthInfo)
	assert.Equal(t, ctx.Cluster, got.Cluster)
	assert.Equal(t, ctx.Namespace, got.Namespace)

	ctx.Namespace = ""
	err = cfg.SetContext(ctxname, ctx)
	assert.NoError(t, err)

	got, err = cfg.GetContext(ctxname)
	assert.NoError(t, err)
	assert.Equal(t, "", got.Namespace)

	err = cfg.RemoveContext(ctxname)
	assert.NoError(t, err)

	_, err = cfg.GetContext(ctxname)
	assert.Contains(t, err.Error(), "not found")
}

func TestConfigSetCluster(t *testing.T) {
	c := "testdata/config"
	clustername := "testset"
	cluster := &clientcmdapi.Cluster{
		Server:                "https://7.8.9.0",
		InsecureSkipTLSVerify: false,
		CertificateAuthority:  "testset-fake-ca-file",
	}
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &c
	cfg := NewConfig(flags)

	err := cfg.SetCluster(clustername, cluster)
	assert.NoError(t, err)

	got, err := cfg.GetCluster(clustername)
	assert.NoError(t, err)

	assert.Equal(t, cluster.Server, got.Server)
	assert.Equal(t, cluster.InsecureSkipTLSVerify, got.InsecureSkipTLSVerify)
	assert.Contains(t, got.CertificateAuthority, cluster.CertificateAuthority)

	cluster.Server = "https://10.9.8.7"

	err = cfg.SetCluster(clustername, cluster)
	assert.NoError(t, err)

	got, err = cfg.GetCluster(clustername)
	assert.NoError(t, err)
	assert.Equal(t, "https://10.9.8.7", got.Server)

	err = cfg.RemoveCluster(clustername)
	assert.NoError(t, err)

	_, err = cfg.GetContext(clustername)
	assert.Contains(t, err.Error(), "not found")
}

func TestConfigSetCredentials(t *testing.T) {
	c := "testdata/config"
	authname := "testset"
	auth := &clientcmdapi.AuthInfo{
		ClientCertificateData: []byte("==ezzdifugh"),
		ClientKeyData:         []byte("ezzdifugh=="),
	}
	flags := genericclioptions.NewConfigFlags(false)
	flags.KubeConfig = &c
	cfg := NewConfig(flags)

	err := cfg.SetCredential(authname, auth)
	assert.NoError(t, err)

	got, err := cfg.GetCredential(authname)
	assert.NoError(t, err)

	assert.Equal(t, auth.ClientCertificate, got.ClientCertificate)
	assert.Equal(t, auth.ClientKeyData, got.ClientKeyData)

	auth.ClientKeyData = []byte("ezzdifugh==ezzdifugh==")

	err = cfg.SetCredential(authname, auth)
	assert.NoError(t, err)

	got, err = cfg.GetCredential(authname)
	assert.NoError(t, err)
	assert.Equal(t, "ezzdifugh==ezzdifugh==", string(got.ClientKeyData))

	err = cfg.RemoveCredential(authname)
	assert.NoError(t, err)

	_, err = cfg.GetCredential(authname)
	assert.Contains(t, err.Error(), "not found")
}
